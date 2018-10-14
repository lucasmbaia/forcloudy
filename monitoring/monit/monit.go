package monit

import (
	"context"
	"encoding/json"
	"forcloudy/monitoring/docker"
	"forcloudy/monitoring/kafka"
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

type Config struct {
	Running int
	Topic   string
	Key     string
	Kafka	[]string
}

func Run(ctx context.Context, conf Config) error {
	var (
		containers  []docker.Containers
		err         error
		errc        = make(chan error, 1)
		t           *time.Ticker
		hostname    string
		metrics     _metrics.Metrics
		mcontainers []_metrics.Containers
		net         = utils.NewNetwork()
		wg          sync.WaitGroup
		body        []byte
		producer    *kafka.Producer
		message     = make(chan []byte)
	)

	if hostname, err = os.Hostname(); err != nil {
		errc <- err
	}

	if producer, err = kafka.NewProducer(context.Background(), conf.Kafka, 5); err != nil {
		errc <- err
	}

	if containers, err = docker.ListAllContainers(NO_FILTER); err != nil {
		errc <- err
	}

	go func() {
		if err = producer.SyncProducer(conf.Topic, conf.Key, message); err != nil {
			errc <- err
		}
	}()

	go updateContainers(ctx, &containers, net)
	for _, container := range containers {
		if err = net.SetPid(container.PID); err != nil {
			errc <- err
		}
	}

	t = time.NewTicker(time.Duration(conf.Running) * time.Second)
	metrics.Hostname = hostname

	go func() {
		for {
			select {
			case _ = <-t.C:
				mcontainers = []_metrics.Containers{}
				wg.Add(len(containers))

				for _, container := range containers {
					go func(container docker.Containers) {
						var (
							wgContainer sync.WaitGroup
							networks    []_metrics.Networks
							cpu         _metrics.Cpu
							memory      _metrics.Memory
						)

						wgContainer.Add(3)

						go func() {
							if networks, err = net.NetworkUtilization(container.PID, conf.Running); err != nil {
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
						mcontainers = append(mcontainers, _metrics.Containers{
							ID:       container.ID,
							Name:     container.Name,
							Cpu:      cpu,
							Memory:   memory,
							Networks: networks,
						})

						wg.Done()
					}(container)
				}

				wg.Wait()
				metrics.Customers = listCustomers(mcontainers)

				if body, err = json.Marshal(mcontainers); err != nil {
					log.Println(err)
					break
				}

				for _, customer := range metrics.Customers {
					if body, err = json.Marshal(customer); err != nil {
						log.Println(err)
						continue
					}

					log.Println(string(body))
					message <- body
				}

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

func listCustomers(containers []_metrics.Containers) []_metrics.Customers {
	var (
		cn        []string
		as        []string
		app       string
		containsC bool
		containsA bool
		customers []_metrics.Customers
	)

	for _, container := range containers {
		cn = strings.Split(container.Name, "_app-")
		as = strings.Split(cn[1], "-")
		app = strings.Join(as[:len(as)-1], "-")
		containsA = false
		containsC = false

		for indexC, customer := range customers {
			if customer.Name == cn[0] {
				for indexA, application := range customers[indexC].Applications {
					if app == application.Name {
						customers[indexC].Applications[indexA].Containers = append(customers[indexC].Applications[indexA].Containers, container)

						containsA = true
						break
					}
				}

				if !containsA {
					customers[indexC].Applications = append(customers[indexC].Applications, _metrics.Applications{
						Name:       app,
						Containers: []_metrics.Containers{container},
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
					Name:       app,
					Containers: []_metrics.Containers{container},
				}},
			})
		}
	}

	return customers
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
