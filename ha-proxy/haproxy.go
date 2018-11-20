package haproxy

import (
	"fmt"
	"github.com/lucasmbaia/forcloudy/api/config"
)

const (
	KEY_ETCD = "/fc-haproxy/"
)

var (
	PROTOCOL_HTTP = []string{"http", "https"}
	PORTS_HTTP    = []string{"80", "443"}
	KEY_HTTP      = map[string]string{"80": "app-http", "443": "app-https"}
)

type Haproxy struct {
	Customer         string
	ApplicationName  string
	ContainerName    string
	PortsContainer   map[string][]string
	Protocol         map[string]string
	AddressContainer string
	Dns              string
	Minion           string
}

type infos struct {
	Customer          string
	ApplicationName   string
	ContainerName     string
	PortSource        string
	PortsDestionation []string
	AddressContainer  string
	Dns               string
	Protocol          string
	Minion            string
}

type ConfHttpHttps struct {
	Hosts []Hosts `json:"hosts,omitempty"`
}

type ConfTcpUdp struct {
	Hosts []Hosts `json:"hosts,omitempty"`
	Dns   string  `json:"dns,omitempty"`
}

type Hosts struct {
	Containers []Containers `json:"containers,omitempty"`
	Name       string       `json:"name,omitempty"`
	Dns        string       `json:"dns,omitempty"`
	Minion     string       `json:"minion,omitempty"`
	PortSRC    string       `json:"portSRC,omitempty"`
	Protocol   string       `json:"protocol,omitempty"`
}

type Containers struct {
	Minion  string `json:"minion,omitempty"`
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

func GenerateConf(h Haproxy) error {
	var (
		exists bool
		key    string
		err    error
	)

	for src, dst := range h.PortsContainer {
		if _, exists = ExistsStringElement(src, PORTS_HTTP); exists {
			var confHttpHttps ConfHttpHttps

			if confHttpHttps, err = httpAndHttps(infos{
				ApplicationName:   h.ApplicationName,
				ContainerName:     h.ContainerName,
				PortSource:        src,
				PortsDestionation: dst,
				AddressContainer:  h.AddressContainer,
				Dns:               h.Dns,
				Minion:            h.Minion,
			}); err != nil {
				return err
			}

			key = fmt.Sprintf("%s%s", KEY_ETCD, KEY_HTTP[src])
			if err = config.EnvSingleton.EtcdConnection.Set(key, confHttpHttps); err != nil {
				return err
			}
		} else {
			var confTcpUdp ConfTcpUdp

			if confTcpUdp, err = tcpAndUdp(infos{
				Customer:          h.Customer,
				ApplicationName:   h.ApplicationName,
				ContainerName:     h.ContainerName,
				PortSource:        src,
				PortsDestionation: dst,
				Dns:               h.Dns,
				Protocol:          h.Protocol[src],
				Minion:            h.Minion,
			}); err != nil {
				return err
			}

			key = fmt.Sprintf("%s%s/%s", KEY_ETCD, h.Customer, h.ApplicationName)
			if err = config.EnvSingleton.EtcdConnection.Set(key, confTcpUdp); err != nil {
				return err
			}
		}
	}

	return nil
}

func RemoveContainer(h Haproxy) error {
	var (
		exists bool
		key    string
		err    error
	)

	for src, _ := range h.PortsContainer {
		if _, exists = ExistsStringElement(src, PORTS_HTTP); exists {
			key = fmt.Sprintf("%s%s", KEY_ETCD, KEY_HTTP[src])

			if exists = config.EnvSingleton.EtcdConnection.Exists(key); exists {
				var conf ConfHttpHttps

				if err = config.EnvSingleton.EtcdConnection.Get(key, &conf); err != nil {
					return err
				}

				for idxHost, host := range conf.Hosts {
					if h.ApplicationName == host.Name {
						for idxContainer, container := range host.Containers {
							if container.Name == h.ContainerName {
								if len(host.Containers)-1 == idxContainer {
									conf.Hosts[idxHost].Containers = conf.Hosts[idxHost].Containers[:idxContainer]
								} else {
									conf.Hosts[idxHost].Containers = append(conf.Hosts[idxHost].Containers[:idxContainer], conf.Hosts[idxHost].Containers[idxContainer+1:]...)
								}
							}
						}
					}
				}

				if err = config.EnvSingleton.EtcdConnection.Set(key, conf); err != nil {
					return err
				}
			}
		} else {
			key = fmt.Sprintf("%s%s/%s", KEY_ETCD, h.Customer, h.ApplicationName)

			if exists = config.EnvSingleton.EtcdConnection.Exists(key); exists {
				var conf ConfHttpHttps

				if err = config.EnvSingleton.EtcdConnection.Get(key, &conf); err != nil {
					return err
				}

				for idxHost, host := range conf.Hosts {
					if host.PortSRC == src {
						for idxContainer, container := range host.Containers {
							if container.Name == h.ContainerName {
								if len(host.Containers)-1 == idxContainer {
									conf.Hosts[idxHost].Containers = conf.Hosts[idxHost].Containers[:idxContainer]
								} else {
									conf.Hosts[idxHost].Containers = append(conf.Hosts[idxHost].Containers[:idxContainer], conf.Hosts[idxHost].Containers[idxContainer+1:]...)
								}
							}
						}
					}
				}

				if err = config.EnvSingleton.EtcdConnection.Set(key, conf); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func tcpAndUdp(i infos) (ConfTcpUdp, error) {
	var (
		key      string
		exists   bool
		conf     ConfTcpUdp
		contains bool
		err      error
	)

	key = fmt.Sprintf("%s%s/%s", KEY_ETCD, i.Customer, i.ApplicationName)
	exists = config.EnvSingleton.EtcdConnection.Exists(key)

	if exists {
		if err = config.EnvSingleton.EtcdConnection.Get(key, &conf); err != nil {
			return conf, err
		}

		for idx, host := range conf.Hosts {
			if host.PortSRC == i.PortSource {
				contains = true
				for _, port := range i.PortsDestionation {
					conf.Hosts[idx].Containers = append(conf.Hosts[idx].Containers, Containers{
						Name:    i.ContainerName,
						Address: fmt.Sprintf("%s:%s", i.AddressContainer, port),
						Minion:  i.Minion,
					})
				}
			}
		}

		if !contains {
			var containers []Containers
			for _, port := range i.PortsDestionation {
				containers = append(containers, Containers{
					Name:    i.ContainerName,
					Address: fmt.Sprintf("%s:%s", i.AddressContainer, port),
					Minion:  i.Minion,
				})
			}

			conf.Hosts = append(conf.Hosts, Hosts{
				PortSRC:    i.PortSource,
				Protocol:   i.Protocol,
				Dns:        i.Dns,
				Containers: containers,
			})
		}
	} else {
		conf = ConfTcpUdp{
			Dns: i.Dns,
			Hosts: []Hosts{
				{PortSRC: i.PortSource, Protocol: i.Protocol},
			},
		}

		for _, port := range i.PortsDestionation {
			conf.Hosts[0].Containers = append(conf.Hosts[0].Containers, Containers{
				Name:    i.ContainerName,
				Address: fmt.Sprintf("%s:%s", i.AddressContainer, port),
				Minion:  i.Minion,
			})
		}
	}

	return conf, nil
}

func httpAndHttps(h infos) (ConfHttpHttps, error) {
	var (
		key      string
		exists   bool
		conf     ConfHttpHttps
		contains bool
		err      error
	)

	key = fmt.Sprintf("%s%s", KEY_ETCD, KEY_HTTP[h.PortSource])
	exists = config.EnvSingleton.EtcdConnection.Exists(key)

	if exists {
		if err = config.EnvSingleton.EtcdConnection.Get(key, &conf); err != nil {
			return conf, err
		}

		for idx, host := range conf.Hosts {
			if host.Name == h.ApplicationName {
				contains = true
				for _, port := range h.PortsDestionation {
					conf.Hosts[idx].Containers = append(conf.Hosts[idx].Containers, Containers{
						Minion:  h.Minion,
						Name:    h.ContainerName,
						Address: fmt.Sprintf("%s:%s", h.AddressContainer, port),
					})
				}
			}
		}

		if !contains {
			var containers []Containers
			for _, port := range h.PortsDestionation {
				containers = append(containers, Containers{
					Minion:  h.Minion,
					Name:    h.ContainerName,
					Address: fmt.Sprintf("%s:%s", h.AddressContainer, port),
				})
			}

			conf.Hosts = append(conf.Hosts, Hosts{
				Name:       h.ApplicationName,
				Dns:        h.Dns,
				Minion:     h.Minion,
				Containers: containers,
			})
		}
	} else {
		conf = ConfHttpHttps{
			Hosts: []Hosts{
				{Name: h.ApplicationName, Dns: h.Dns, Minion: h.Minion},
			},
		}

		for _, port := range h.PortsDestionation {
			conf.Hosts[0].Containers = append(conf.Hosts[0].Containers, Containers{Minion: h.Minion, Name: h.ContainerName, Address: fmt.Sprintf("%s:%s", h.AddressContainer, port)})
		}
	}

	return conf, nil
}

func ExistsStringElement(f string, s []string) (int, bool) {
	for idx, str := range s {
		if str == f {
			return idx, true
		}
	}

	return 0, false
}
