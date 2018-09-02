package stats

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type Metrics struct {
}

func stats(update int, metrics <-chan Metrics) error {
	var (
		output     []byte
		err        error
		timer      <-chan time.Time
		command    *exec.Cmd
		containers []string
	)

	timer = time.Tick(time.Duration(update) * time.Second)

	for now := range timer {
		var (
			infos        []string
			memory       []string
			cn           []string
			as           []string
			app          string
			cCustomer    bool
			cApplication bool
		)

		command = exec.Command("docker", "stats", "--all", "--no-stream", "--format", "'{{.Container}}#{{ .Name }}#{{.CPUPerc}}#{{.MemUsage}}'")
		if output, err = command.CombinedOutput(); err != nil {
			return err
		}

		containers = strings.Split(strings.Replace(string(output), "'", "", -1), "\n")

		for _, container := range containers {
			if len(container) > 0 {
				infos = strings.Split(container, "#")
				cCustomer = false
				cApplication = false

				memory = strings.Split(infos[3], "/")
				cn = strings.Split(strings.TrimSpace(infos[1]), "_app-")
				as = strings.Split(cn[1], "-")
				app = strings.Join(as[:len(as)-1], "-")

				fmt.Println(infos, memory, cn, as, app, cCustomer, cApplication)
			}
		}

		fmt.Println(now)
	}

	return nil
}

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

func ListAllContainers() ([]Containers, error) {
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
