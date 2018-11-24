package core

import (
	"fmt"
	dockerxmpp "github.com/lucasmbaia/go-xmpp/docker"
	"testing"
)

func TestNameContainers(t *testing.T) {
	var (
		elements dockerxmpp.Elements
		err      error
	)

	if elements, err = nameContainers(); err != nil {
		t.Fatal(err)
	}

	fmt.Println(elements)
}

func TestTotalContainers(t *testing.T) {
	var (
		elements dockerxmpp.Elements
		err      error
	)

	if elements, err = totalContainers(); err != nil {
		t.Fatal(err)
	}

	fmt.Println(elements)
}

func TestExistsImage(t *testing.T) {
	var (
		elements dockerxmpp.Elements
		err      error
	)

	elements.Name = "openstack-netpartition/image:v1"

	if elements, err = existsImage(elements); err != nil {
		t.Fatal(err)
	}

	fmt.Println(elements)
}

func TestLoadImage(t *testing.T) {
	var (
		elements dockerxmpp.Elements
		err      error
	)

	elements.Name = "lucas.tar.gz"

	if _, err = loadImage(elements); err != nil {
		t.Fatal(err)
	}
}

func TestGenerateImage(t *testing.T) {
	var (
		elements dockerxmpp.Elements
		err      error
	)

	elements.Name = "lucas_app-bematech"
	elements.BuildName = "hello_world"
	elements.Tag = "v1"

	if err = generateImage(elements); err != nil {
		t.Fatal(err)
	}
}

func TestDeploy(t *testing.T) {
	var (
		elements dockerxmpp.Elements
		err      error
	)

	elements = dockerxmpp.Elements{
		Name:   "lucas_app-bematech",
		Cpus:   "0.1",
		Memory: "15MB",
		Image:  "alpine",
		Ports: []dockerxmpp.Ports{
			{Port: 80},
		},
	}

	if _, err = deploy(elements, true); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveContainer(t *testing.T) {
	if err := removeContainer(dockerxmpp.Elements{Name: "bematech"}); err != nil {
		t.Fatal(err)
	}
}

func TestAddressContainer(t *testing.T) {
	if address, err := addressContainer(dockerxmpp.Elements{Name: "lucas"}); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(address.Address)
	}
}

func TestPortsContainer(t *testing.T) {
	if ports, err := portsContainer(dockerxmpp.Elements{Name: "lucas-2"}); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(ports.PortsContainer)
	}
}

func TestExistsContainer(t *testing.T) {
	if exists, err := existsContainer(dockerxmpp.Elements{Name: "lucas-2"}); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(exists)
	}
}
