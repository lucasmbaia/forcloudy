package core

import (
	"fmt"
	"strings"

	"github.com/lucasmbaia/forcloudy/minion/utils"
	dockerxmpp "github.com/lucasmbaia/go-xmpp/docker"
)

const (
	TIMEOUT_DEFAULT_COMMAND = 60
)

var (
	stepsGenerateImage = []string{
		"docker exec -it {container} mkdir /app",
		"docker exec -it {container} apk add --no-cache bash",
		"docker cp {application} {container}:/app",
		"docker commit --change='ENTRYPOINT [\"/app/{application}\"]' {container} {container}/image:{tag}",
		"docker save {container}/image:{tag} -o /images/{container}.tar.gz",
	}
)

func generateImage(elements dockerxmpp.Elements) error {
	var (
		err error
	)

	if _, err = utils.Command("docker", []string{"exec", "-t", elements.Name, "mkdir", "/app"}, TIMEOUT_DEFAULT_COMMAND); err != nil {
		return err
	}

	if _, err = utils.Command("docker", []string{"exec", "-t", elements.Name, "apk", "add", "--no-cache", "bash"}, 120); err != nil {
		return err
	}

	if _, err = utils.Command("docker", []string{"cp", fmt.Sprintf("%s%s", elements.Path, elements.BuildName), fmt.Sprintf("%s:/app", elements.Name)}, TIMEOUT_DEFAULT_COMMAND); err != nil {
		return err
	}

	if _, err = utils.Command("docker", []string{"commit", "--change", fmt.Sprintf("ENTRYPOINT [\"/app/%s\"]", elements.BuildName), elements.Name, fmt.Sprintf("%s/image:%s", elements.Name, elements.Tag)}, TIMEOUT_DEFAULT_COMMAND); err != nil {
		return err
	}

	if _, err = utils.Command("docker", []string{"save", fmt.Sprintf("%s/image:%s", elements.Name, elements.Tag), "-o", fmt.Sprintf("/images/%s.tar.gz", elements.Name)}, TIMEOUT_DEFAULT_COMMAND); err != nil {
		return err
	}

	return nil
}

func loadImage(elements dockerxmpp.Elements) (dockerxmpp.Elements, error) {
	var (
		err   error
		image string
	)

	if elements.Path != "" {
		if elements.Path[len(elements.Path)-1:] != "/" {
			elements.Path = fmt.Sprintf("%s/", elements.Path)
		}
	}

	image = fmt.Sprintf("%s%s", elements.Path, elements.Name)

	if _, err = utils.Command("docker", []string{"load", "--input", image}, TIMEOUT_DEFAULT_COMMAND); err != nil {
		return dockerxmpp.Elements{}, err
	}

	return dockerxmpp.Elements{}, nil
}

func existsImage(elements dockerxmpp.Elements) (dockerxmpp.Elements, error) {
	var (
		result []string
		err    error
		name   string
	)

	if result, err = utils.Command("docker", []string{"images", "--format", "{{.Repository}}:{{.Tag}}"}, TIMEOUT_DEFAULT_COMMAND); err != nil {
		return elements, err
	}

	if len(strings.Split(elements.Name, ":")) == 1 {
		name = fmt.Sprintf("%s:latest", name)
	} else {
		name = elements.Name
	}

	for _, image := range result {
		if image == name {
			return elements, nil
		}
	}

	elements.Name = EMPTY_STR
	return elements, nil
}

func deploy(elements dockerxmpp.Elements, imageCreate bool) (dockerxmpp.Elements, error) {
	var (
		result []string
		err    error
		args   []string
	)

	if imageCreate {
		args = []string{"run", "-t", "--rm"}
	} else {
		args = []string{"run", "--rm"}
	}

	if len(elements.Args) > 0 {
		for _, arg := range elements.Args {
			args = append(args, "--env")
			args = append(args, fmt.Sprintf("%s=%s", arg.Name, arg.Value))
		}
	}

	if len(elements.Ports) > 0 {
		args = append(args, "-P")
		for _, port := range elements.Ports {
			args = append(args, fmt.Sprintf("--expose=%d", port.Port))
		}
	}

	args = append(args, []string{"--name", elements.Name}...)
	args = append(args, fmt.Sprintf("--cpus=%s", elements.Cpus))
	args = append(args, fmt.Sprintf("--memory=%s", elements.Memory))
	args = append(args, []string{"-d", elements.Image}...)

	if result, err = utils.Command("docker", args, TIMEOUT_DEFAULT_COMMAND); err != nil {
		return dockerxmpp.Elements{}, err
	}

	return dockerxmpp.Elements{ID: result[0], Name: elements.Name}, nil
}

func masterDeploy(elements dockerxmpp.Elements) (dockerxmpp.Elements, error) {
	return deploy(elements, elements.CreateImage)
}

func appendDeploy(elements dockerxmpp.Elements) (dockerxmpp.Elements, error) {
	return deploy(elements, false)
}

func nameContainers() (dockerxmpp.Elements, error) {
	var (
		elements dockerxmpp.Elements
		result   []string
		err      error
	)

	if result, err = utils.Command("docker", []string{"ps", "-a", "--format", "{{.Names}}"}, TIMEOUT_DEFAULT_COMMAND); err != nil {
		return elements, err
	}

	for _, container := range result {
		elements.Containers = append(elements.Containers, dockerxmpp.Container{Name: container})
	}

	return elements, nil
}

func totalContainers() (dockerxmpp.Elements, error) {
	var (
		elements dockerxmpp.Elements
		result   []string
		err      error
	)

	if result, err = utils.Command("docker", []string{"ps", "-a", "--format", "{{.Names}}"}, TIMEOUT_DEFAULT_COMMAND); err != nil {
		return elements, err
	}

	elements.TotalContainers = len(result)

	return elements, nil
}

func operationContainers(elements dockerxmpp.Elements) (dockerxmpp.Elements, error) {
	return dockerxmpp.Elements{}, nil
}

func removeContainer(elements dockerxmpp.Elements) error {
	var (
		err    error
		exists bool
	)

	if exists, err = existsContainer(elements); err != nil {
		return err
	}

	if exists {
		if _, err = utils.Command("docker", []string{"rm", "-f", elements.Name}, TIMEOUT_DEFAULT_COMMAND); err != nil {
			return err
		}
	}

	return nil
}

func existsContainer(elements dockerxmpp.Elements) (bool, error) {
	var (
		err error
		el  dockerxmpp.Elements
	)

	if el, err = nameContainers(); err != nil {
		return false, err
	}

	for _, container := range el.Containers {
		if container.Name == elements.Name {
			return true, nil
		}
	}

	return false, nil
}

func addressContainer(elements dockerxmpp.Elements) (dockerxmpp.Elements, error) {
	var (
		err    error
		result []string
	)

	if result, err = utils.Command("docker", []string{"inspect", "-f", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}", elements.Name}, TIMEOUT_DEFAULT_COMMAND); err != nil {
		return dockerxmpp.Elements{}, err
	}

	return dockerxmpp.Elements{Address: result[0]}, nil
}

func portsContainer(elements dockerxmpp.Elements) (dockerxmpp.Elements, error) {
	var (
		result []string
		err    error
		ports  []string
		pc     = make(map[string][]string)
	)

	if result, err = utils.Command("docker", []string{"inspect", "--format={{range $p, $conf := .NetworkSettings.Ports}}{{$p}}:{{(index $conf 0).HostPort}}-{{end}}", elements.Name}, TIMEOUT_DEFAULT_COMMAND); err != nil {
		return dockerxmpp.Elements{}, err
	}

	ports = strings.Split(result[0], "-")
	ports = ports[:len(ports)-1]

	for _, port := range ports {
		var p = strings.Split(port, "/tcp:")

		if _, ok := pc[p[0]]; ok {
			pc[p[0]] = append(pc[p[0]], p[1])
		} else {
			pc[p[0]] = []string{p[1]}
		}
	}

	for src, dst := range pc {
		elements.PortsContainer = append(elements.PortsContainer, dockerxmpp.PortsContainer{Source: src, Destinations: dst})
	}

	return elements, nil
}
