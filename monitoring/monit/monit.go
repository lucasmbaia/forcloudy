package monit

import (
	"context"
	"encoding/json"
	"forcloudy/monitoring/docker"
	_metrics "forcloudy/monitoring/metrics"
	"forcloudy/monitoring/utils"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	NO_FILTER = ""
	START     = "start"
	DIE       = "die"
)

func Run(ctx context.Context, running int) error {
	var (
		containers []docker.Containers
		err        error
		errc       = make(chan error, 1)
		t          *time.Ticker
		hostname   string
		metrics    _metrics.Metrics
		customers  []_metrics.Customers
		net        = utils.NewNetwork()
		wg         sync.WaitGroup
		body       []byte
	)

	if hostname, err = os.Hostname(); err != nil {
		errc <- err
	}

	if containers, err = docker.ListAllContainers(NO_FILTER); err != nil {
		errc <- err
	}

	go updateContainers(ctx, &containers, net)
	for _, container := range containers {
		if err = net.SetPid(container.PID); err != nil {
			errc <- err
		}
	}

	t = time.NewTicker(time.Duration(running) * time.Second)
	metrics.Hostname = hostname

	go func() {
		for {
			select {
			case _ = <-t.C:
				customers = []_metrics.Customers{}
				wg.Add(len(containers))

				for _, container := range containers {
					go func(container docker.Containers) {
						var (
							cn          []string
							as          []string
							app         string
							containsC   bool
							containsA   bool
							wgContainer sync.WaitGroup
							networks    []_metrics.Networks
							cpu         _metrics.Cpu
							memory      _metrics.Memory
						)

						wgContainer.Add(3)
						cn = strings.Split(container.Name, "_app-")
						as = strings.Split(cn[1], "-")
						app = strings.Join(as[:len(as)-1], "-")

						go func() {
							if networks, err = net.NetworkUtilization(container.PID, running); err != nil {
								log.Println(err)
							}
							wgContainer.Done()
						}()

						go func() {
							if memory, err = utils.MemoryUtilization(container.ID); err != nil {
								log.Println(err)
							}
							wgContainer.Done()
						}()

						go func() {
							if cpu.TotalUsage, err = utils.CpuUsageContainerUnix(container.ID, 1); err != nil {
								log.Println(err)
							}
							wgContainer.Done()
						}()

						wgContainer.Wait()
						for indexC, customer := range customers {
							if customer.Name == cn[0] {
								for indexA, application := range customers[indexC].Applications {
									if app == application.Name {
										customers[indexC].Applications[indexA].Containers = append(customers[indexC].Applications[indexA].Containers, _metrics.Containers{
											ID:       container.ID,
											Name:     container.Name,
											Cpu:      cpu,
											Memory:   memory,
											Networks: networks,
										})

										containsA = true
										break
									}
								}

								if !containsA {
									customers[indexC].Applications = append(customers[indexC].Applications, _metrics.Applications{
										Name: app,
										Containers: []_metrics.Containers{{
											ID:       container.ID,
											Name:     container.Name,
											Cpu:      cpu,
											Memory:   memory,
											Networks: networks,
										}},
									})
								}

								containsC = true
								break
							}
						}

						if !containsC {
							customers = append(customers, _metrics.Customers{
								Name: cn[0],
								Applications: []_metrics.Applications{{
									Name: app,
									Containers: []_metrics.Containers{{
										ID:       container.ID,
										Name:     container.Name,
										Cpu:      cpu,
										Memory:   memory,
										Networks: networks,
									}},
								}},
							})
						}

						wg.Done()
					}(container)
				}

				wg.Wait()
				metrics.Customers = customers

				if body, err = json.Marshal(metrics); err != nil {
					log.Println(err)
					break
				}

				log.Println(string(body))
			case _ = <-ctx.Done():
				log.Println("DONE")
			}
		}
	}()

	select {
	case err = <-errc:
		return err
	case _ = <-ctx.Done():
		return nil
	}
}

func updateContainers(ctx context.Context, containers *[]docker.Containers, net *utils.Utilization) {
	var (
		err   error
		event = make(chan []byte)
		errc  = make(chan error, 1)
		c     []docker.Containers
	)

	go func() {
		errc <- docker.DockerEvents(ctx, event)
	}()

	for {
		select {
		case msg := <-event:
			var ev docker.Events

			if err = json.Unmarshal(msg, &ev); err != nil {
				log.Println(err)
				continue
			}

			switch ev.Status {
			case START:
				if c, err = docker.ListAllContainers(ev.Actor.Attributes.Name); err != nil {
					log.Println(err)
					break
				}

				if err = net.SetPid(c[0].PID); err != nil {
					log.Println(err)
					break
				}

				*containers = append(*containers, c...)
			case DIE:
				for index, container := range *containers {
					if container.Name == ev.Actor.Attributes.Name {
						*containers = append((*containers)[:index], (*containers)[index+1:]...)
						net.DelPid(container.PID)
						break
					}
				}
			}
		case e := <-errc:
			log.Println(e)
		}
	}
}
