package core

import (
	"fmt"
	"github.com/lucasmbaia/forcloudy/api/config"
	"github.com/lucasmbaia/go-xmpp/docker"
	"sync"
)

const (
	IMAGE_BASE  = "alpine"
	PATH_IMAGES = "/images/"
)

type Deploy struct {
	Customer        string            `json:",omitempty"`
	ApplicationName string            `json:",omitempty"`
	ImageVersion    string            `json:",omitempty"`
	Ports           []Ports           `json:",omitempty"`
	Args            map[string]string `json:",omitempty"`
	Cpus            string            `json:",omitempty"`
	Memory          string            `json:",omitempty"`
	TotalContainers int               `json:",omitempty"`
	Dns             string            `json:",omitempty"`
	Image           string            `json:",omitempty"`
	Build           string            `json:",omitempty"`
	Path            string            `json:",omitempty"`
}

type Ports struct {
	Port     int    `json:",omitempty"`
	Protocol string `json:",omitempty"`
}

type MinionsCount struct {
	Name            string
	TotalContainers int
}

type Container struct {
	ID             string
	Name           string
	Address        string
	PortsContainer []docker.PortsContainer
	Minion         string
	Error          error
}

func DeployApplication(d Deploy, iterator int, first bool, assyncContainers chan<- Container) error {
	var (
		image           string
		applicationName string
		exists          bool
		gImage          = true
		err             error
		uploadImage     []string
		listMinions     []string
		wg              sync.WaitGroup
		errc            = make(chan error, 1)
		minionsCount    map[string]int
	)

	image = fmt.Sprintf("%s_app-%s/image:%s", d.Customer, d.ApplicationName, d.ImageVersion)
	applicationName = fmt.Sprintf("%s_app-%s", d.Customer, d.ApplicationName)
	d.Image = image

	if first {
		for minion, _ := range minions {
			listMinions = append(listMinions, minion)
			if exists, err = existsImage(image, minion); err != nil {
				return err
			}

			if exists {
				gImage = false
			} else {
				uploadImage = append(uploadImage, minion)
			}
		}

		if gImage {
			if _, err = createContainer(d, listMinions[0], applicationName, IMAGE_BASE, true); err != nil {
				return err
			}

			if err = generateImage(applicationName, d.ImageVersion, d.Path, d.Build, listMinions[0]); err != nil {
				return err
			}

			if err = removeContainer(applicationName, listMinions[0]); err != nil {
				return err
			}

			listMinions = listMinions[1:]
		}

		for _, minion := range uploadImage {
			if err = loadImage(PATH_IMAGES, fmt.Sprintf("%s.tar.gz", applicationName), minion); err != nil {
				return err
			}
		}
	}

	if len(minions) == 1 {
		var minion string
		for m, _ := range minions {
			minion = m
		}

		wg.Add(d.TotalContainers)

		for i := iterator; i <= d.TotalContainers; i++ {
			var containerName = fmt.Sprintf("%s_app-%s-%d", d.Customer, d.ApplicationName, i)
			go func(containerName string) {
				var c Container
				if c, err = createContainer(d, minion, containerName, image, false); err != nil {
					c.Error = err
					errc <- err
				}

				if assyncContainers != nil {
					assyncContainers <- c
				}
				wg.Done()
			}(containerName)
		}

		go func() {
			select {
			case e := <-errc:
				err = e
			}
		}()

		wg.Wait()
	} else {
		if minionsCount, err = containersPerMinion(d.TotalContainers); err != nil {
			return err
		}

		wg.Add(d.TotalContainers)

		for key, value := range minionsCount {
			for i := 0; i < value; i++ {
				var containerName = fmt.Sprintf("%s_app-%s-%d", d.Customer, d.ApplicationName, iterator)
				go func(containerName, minion string) {
					var c Container
					if c, err = createContainer(d, minion, containerName, image, false); err != nil {
						c.Error = err
						errc <- err
					}

					if assyncContainers != nil {
						assyncContainers <- c
					}
					wg.Done()
				}(containerName, key)

				iterator++
			}
		}

		go func() {
			select {
			case e := <-errc:
				err = e
			}
		}()

		wg.Wait()
	}

	return nil
}

func RemoveContainer(name string) error {
	var (
		err    error
		exists bool
		idx    int
	)

	for m, i := range minions {
		if idx, exists = utils.ExistsStringElement(name, i.Containers); exists {
			if err = removeContainer(name, m); err != nil {
				return err
			}

			if len(i.Containers)-1 == idx {
				i.Containers = i.Containers[:idx]
			} else {
				i.Containers = append(i.Containers[:idx], i.Containers[idx+1:]...)
			}

			minions[m] = Minions{Containers: i}
			break
		}
	}

	return nil
}

func existsImage(image, to string) (bool, error) {
	var (
		iq       docker.IQ
		err      error
		response = make(chan ResponseIQ, 1)
	)

	if iq, err = docker.ExistsImage(docker.Image{
		From: config.EnvSingleton.XmppConnection.Jid,
		To:   to,
		Name: image,
	}); err != nil {
		return false, err
	}

	mutex.Lock()
	responseIQ[iq.ID] = response
	mutex.Unlock()

	if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
		return false, err
	}

	select {
	case r := <-response:
		if r.Elements.Name != EMPTY_STR {
			return true, nil
		}
	}

	return false, nil
}

func createContainer(d Deploy, to, name, image string, imageCreate bool) (Container, error) {
	var (
		iq       docker.IQ
		err      error
		response = make(chan ResponseIQ, 1)
		ports    []docker.Ports
		args     []docker.Args
	)

	if len(d.Ports) > 0 {
		for _, port := range d.Ports {
			ports = append(ports, docker.Ports{Port: port.Port, Protocol: port.Protocol})
		}
	}

	if len(d.Args) > 0 {
		for key, value := range d.Args {
			args = append(args, docker.Args{Name: key, Value: value})
		}
	}

	if iq, err = docker.MasterDeploy(docker.Deploy{
		From:            config.EnvSingleton.XmppConnection.Jid,
		To:              to,
		Customer:        d.Customer,
		ApplicationName: d.ApplicationName,
		Name:            name,
		Cpus:            d.Cpus,
		Memory:          d.Memory,
		Ports:           ports,
		Args:            args,
		Image:           image,
		CreateImage:     imageCreate,
	}); err != nil {
		return Container{}, err
	}

	mutex.Lock()
	responseIQ[iq.ID] = response
	mutex.Unlock()

	if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
		return Container{}, err
	}

	select {
	case r := <-response:
		return Container{
			ID:             r.Elements.ID,
			Name:           name,
			PortsContainer: r.Elements.PortsContainer,
			Address:        r.Elements.Address,
			Minion:         r.Elements.Minion,
		}, r.Error
	}
}

func generateImage(image, version, path, build, to string) error {
	var (
		iq       docker.IQ
		err      error
		response = make(chan ResponseIQ, 1)
	)

	if iq, err = docker.GenerateImage(docker.Image{
		From:      config.EnvSingleton.XmppConnection.Jid,
		To:        to,
		Name:      image,
		Path:      path,
		BuildName: build,
		Tag:       version,
	}); err != nil {
		return err
	}

	mutex.Lock()
	responseIQ[iq.ID] = response
	mutex.Unlock()

	if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
		return err
	}

	select {
	case _ = <-response:
		return nil
	}
}

func removeContainer(name, to string) error {
	var (
		iq       docker.IQ
		err      error
		response = make(chan ResponseIQ, 1)
	)

	if iq, err = docker.RemoveContainer(docker.Action{
		From:      config.EnvSingleton.XmppConnection.Jid,
		To:        to,
		Container: name,
	}); err != nil {
		return err
	}

	mutex.Lock()
	responseIQ[iq.ID] = response
	mutex.Unlock()

	if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
		return err
	}

	select {
	case r := <-response:
		return r.Error
	}
}

func loadImage(path, name, to string) error {
	var (
		iq       docker.IQ
		err      error
		response = make(chan ResponseIQ, 1)
	)

	if iq, err = docker.LoadImage(docker.Image{
		From: config.EnvSingleton.XmppConnection.Jid,
		To:   to,
		Path: path,
		Name: name,
	}); err != nil {
		return err
	}

	mutex.Lock()
	responseIQ[iq.ID] = response
	mutex.Unlock()

	if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
		return err
	}

	select {
	case r := <-response:
		return r.Error
	}
}

func containersPerMinion(totalContainers int) (map[string]int, error) {
	var (
		minionsCount      = make(map[int]MinionsCount)
		minionsContainers = make(map[string]int)
		err               error
		total             int
		totalMinions      = len(minions)
		count             = 0
		differente        bool
	)

	for minion, _ := range minions {
		if total, err = totalContainersMinion(minion); err != nil {
			return minionsContainers, err
		}

		if len(minionsCount) == 0 {
			minionsCount[count] = MinionsCount{Name: minion, TotalContainers: total}
		} else {
			for key, value := range minionsCount {
				if total < value.TotalContainers {
					var mc = MinionsCount{Name: value.Name, TotalContainers: value.TotalContainers}
					minionsCount[key] = MinionsCount{Name: minion, TotalContainers: total}
					minionsCount[count] = mc
				} else {
					minionsCount[count] = MinionsCount{Name: minion, TotalContainers: total}
				}
			}
		}

		count++
	}

	for i := 0; i < totalMinions; i++ {
		for j := i + 1; j < totalMinions; j++ {
			if minionsCount[i].TotalContainers != minionsCount[j].TotalContainers {
				differente = true
				break
			}
		}
	}

	existsMinion := func(n string, m map[string]int) bool {
		for key, _ := range m {
			if key == n {
				return true
			}
		}

		return false
	}

	if !differente {
		for _, value := range minionsCount {
			minionsContainers[value.Name] = totalContainers / totalMinions
		}

		if totalContainers%totalMinions != 0 {
			minionsContainers[minionsCount[0].Name] += totalContainers % totalMinions
		}
	} else {
		for count := 0; count < totalContainers; count++ {
			if count > 0 && minionsCount[0].TotalContainers > minionsCount[1].TotalContainers {
				for idx, _ := range minionsCount {
					x := idx + 1
					if x < len(minionsCount) && minionsCount[idx].TotalContainers > minionsCount[x].TotalContainers {
						total := minionsCount[x].TotalContainers
						name := minionsCount[x].Name
						minionsCount[x] = MinionsCount{Name: minionsCount[idx].Name, TotalContainers: minionsCount[idx].TotalContainers}
						minionsCount[idx] = MinionsCount{Name: name, TotalContainers: total}
					}
				}
			}

			if !existsMinion(minionsCount[0].Name, minionsContainers) {
				minionsContainers[minionsCount[0].Name] = 1
			} else {
				minionsContainers[minionsCount[0].Name] += 1
			}

			minionsCount[0] = MinionsCount{Name: minionsCount[0].Name, TotalContainers: minionsCount[0].TotalContainers + 1}
		}
	}

	return minionsContainers, nil
}

func totalContainersMinion(to string) (int, error) {
	var (
		iq       docker.IQ
		err      error
		response = make(chan ResponseIQ, 1)
	)

	if iq, err = docker.TotalContainers(config.EnvSingleton.XmppConnection.Jid, to); err != nil {
		return 0, err
	}

	mutex.Lock()
	responseIQ[iq.ID] = response
	mutex.Unlock()

	if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
		return 0, err
	}

	select {
	case r := <-response:
		return r.Elements.TotalContainers, r.Error
	}
}
