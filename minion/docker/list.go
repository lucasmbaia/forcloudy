package docker

import (
	"os/exec"
	"strings"
)

type Containers struct {
	ID    string  `json:",omitempty"`
	Name  string  `json:",omitempty"`
	Image string  `json:",omitempty"`
	Ports []Ports `json:",omitempty"`
	PID   string  `json:",omitempty"`
}

type Ports struct {
	Interface   string `json:",omitempty"`
	Source      string `json:",omitempty"`
	Destination string `json:",omitempty"`
}

func ListAllContainers(filter string) ([]Containers, error) {
	var (
		output     []byte
		err        error
		containers []Containers
		ps         []string
		infos      []string
		iAdress    []string
		iPorts     []string
		ds         []string
	)

	if output, err = exec.Command("docker", "ps", "--no-trunc", "--format", "{{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Ports}}").CombinedOutput(); err != nil {
		return containers, err
	}

	ps = strings.Split(string(output), "\n")
	ps = ps[:len(ps)-1]

	for _, container := range ps {
		var ports []Ports
		infos = strings.Split(container, "\t")
		iAdress = strings.Split(strings.TrimSpace(infos[3]), ",")

		if filter != "" {
			if filter != infos[1] {
				continue
			}
		}

		for _, addr := range iAdress {
			if len(addr) > 0 {
				iPorts = strings.Split(addr, ":")
				ds = strings.Split(strings.Replace(iPorts[1], "/tcp", "", -1), "->")
				ports = append(ports, Ports{Interface: iPorts[0], Source: ds[0], Destination: ds[1]})
			}
		}

		if output, err = exec.Command("docker", "inspect", infos[1], "--format", "{{.State.Pid}}").CombinedOutput(); err != nil {
			return containers, err
		}

		containers = append(containers, Containers{ID: infos[0], Name: infos[1], Image: infos[2], Ports: ports, PID: strings.TrimSuffix(string(output), "\n")})
	}

	return containers, nil
}
