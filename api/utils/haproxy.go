package utils

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
	ApplicationName  string
	ContainerName    string
	PortsContainer   map[string][]string
	Protocol         map[string]string
	AddressContainer string
	Dns              string
}

type httpHttps struct {
	ApplicationName   string
	ContainerName     string
	PortSource        string
	PortsDestionation []string
	AddressContainer  string
	Dns               string
	Protocol          string
}

type ConfHttpHttps struct {
	Hosts []Hosts `json:"hosts,omitempty"`
}

type ConfTcpUdp struct {
	Hosts []Hosts `json:"hosts,omitempty"`
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
	/*var (
		exists bool
	)

	for src, dst := range h.PortsContainer {
		if _, exists = ExistsStringElement(src, PORTS_HTTP); exists {

		}
	}*/

	return nil
}

func tcpAndUdp() (ConfTcpUdp, error) {
}

func httpAndHttps(h httpHttps) (ConfHttpHttps, error) {
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
						Minion:  config.EnvConfig.Hostname,
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
					Minion:  config.EnvConfig.Hostname,
					Name:    h.ContainerName,
					Address: fmt.Sprintf("%s:%s", h.AddressContainer, port),
				})
			}

			conf.Hosts = append(conf.Hosts, Hosts{
				Name:       h.ApplicationName,
				Dns:        h.Dns,
				Minion:     config.EnvConfig.Hostname,
				Containers: containers,
			})
		}
	} else {
		conf = ConfHttpHttps{
			Hosts: []Hosts{
				{Name: h.ApplicationName, Dns: h.Dns, Minion: config.EnvConfig.Hostname},
			},
		}

		for _, port := range h.PortsDestionation {
			conf.Hosts[0].Containers = append(conf.Hosts[0].Containers, Containers{Minion: config.EnvConfig.Hostname, Name: h.ContainerName, Address: fmt.Sprintf("%s:%s", h.AddressContainer, port)})
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
