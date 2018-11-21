package core

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/lucasmbaia/forcloudy/dtos"
	"github.com/lucasmbaia/forcloudy/minion/config"
	"github.com/lucasmbaia/forcloudy/minion/docker"
	"github.com/lucasmbaia/forcloudy/minion/log"
	"github.com/lucasmbaia/forcloudy/minion/utils"
	"github.com/lucasmbaia/go-xmpp"
	dockerxmpp "github.com/lucasmbaia/go-xmpp/docker"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	EMPTY_STR           = ""
	JABBER_IQ_DOCKER    = "jabber:iq:docker"
	UNAVAILABLE         = "unavailable"
	AVAILABLE           = "available"
	MAX_RETRY_CONTAINER = 3
)

var (
	masterNode           []string
	containers           []docker.Containers
	containers_deploy    []string
	containers_die       []string
	containersRemove     []string
	eventsDocker         map[string]chan EventsDocker
	containersDie        chan docker.Containers
	retryDeployContainer map[string]int
)

type EventsDocker struct {
	Error error
}

type ApplicationEtcd struct {
	Protocol        map[string]string `json:"protocol,omitempty"`
	Image           string            `json:"image,omitempty"`
	PortsDST        []string          `json:"portsDst,omitempty"`
	Cpus            string            `json:"cpus,omitempty"`
	Dns             string            `json:"dns,omitempty"`
	Memory          string            `json:"memory,omitempty"`
	TotalContainers int               `json:"totalContainers,omitempty"`
}

func init() {
	eventsDocker = make(map[string]chan EventsDocker)
	retryDeployContainer = make(map[string]int)
	containersDie = make(chan docker.Containers)
}

func Run(ctx context.Context) error {
	var (
		err error
	)

	if containers, err = docker.ListAllContainers(EMPTY_STR); err != nil {
		return err
	}

	config.EnvSingleton.Log.Infofc(log.TEMPLATE_ACTION, "Core", "Run", "Containers in this minion", containers)

	//init watch events of docker
	watchEvents(ctx)
	//init check containers die
	checkContainerDie(ctx)

	if err = initXMPP(ctx); err != nil {
		return err
	}

	select {
	case _ = <-ctx.Done():
		return nil
	}
}

func watchEvents(ctx context.Context) {
	var (
		err    error
		errc   = make(chan error, 1)
		event  = make(chan []byte)
		exists bool
	)

	go func() {
		errc <- docker.DockerEvents(context.Background(), event)
	}()

	go func() {
		for {
			select {
			case msg := <-event:
				var ev docker.Events

				if err = json.Unmarshal(msg, &ev); err != nil {
					config.EnvSingleton.Log.Errorf(log.TEMPLATE_ACTION, "Core", "watchEvents", "Unmarshal Event Docker", err.Error())
					break
				}

				if ev.Status == "die" {
					if ev.Actor.Attributes.Name != EMPTY_STR {
						if _, ok := eventsDocker[ev.Actor.Attributes.Name]; ok {
							eventsDocker[ev.Actor.Attributes.Name] <- EventsDocker{Error: errors.New("Erro to create container")}
						} else {
							for _, container := range containers {
								if container.Name == ev.Actor.Attributes.Name {
									if _, exists = utils.ExistsStringElement(ev.Actor.Attributes.Name, containersRemove); exists {
										break
									}

									if container.Image == EMPTY_STR {
										container.Image = ev.Actor.Attributes.Image
									}

									containersDie <- container
									break
								}
							}
						}
					}
				}

				config.EnvSingleton.Log.Debugfc(log.TEMPLATE_ACTION, "Core", "watchEvents", "Event Docker", ev)
			case e := <-errc:
				config.EnvSingleton.Log.Errorf(log.TEMPLATE_ACTION, "Core", "watchEvents", "Receive error of docker events", e.Error())
				watchEvents(ctx)
			}
		}
	}()
}

func checkContainerDie(ctx context.Context) {
	go func() {
		for {
			select {
			case c := <-containersDie:
				if _, ok := retryDeployContainer[c.Name]; !ok {
					retryDeployContainer[c.Name] = 0
				}

				var cn = strings.Split(c.Name, "_app-")
				var customer = cn[0]
				var aux = strings.Split(cn[1], "-")
				var applicationName = strings.Join(aux[:len(aux)-1], "-")
				var key = fmt.Sprintf("/%s/%s", customer, applicationName)
				var ap ApplicationEtcd
				var err error
				var body []byte
				var elements dockerxmpp.Elements
				var message = &xmpp.Message{From: config.EnvSingleton.XmppConnection.Jid}

				if err = config.EnvSingleton.EtcdConnection.Get(key, &ap); err != nil {
					config.EnvSingleton.Log.Errorf(log.TEMPLATE_ACTION, "Core", "checkContainerDie", "Get Infos Etcd", err.Error())
					break
				}

				var container = dtos.Container{
					Customer:    customer,
					Application: applicationName,
					Name:        c.Name,
					Image:       c.Image,
					Ports:       make(map[string][]string),
					Minion:      config.EnvConfig.Hostname,
				}

				if retryDeployContainer[c.Name] < MAX_RETRY_CONTAINER {
					var (
						err   error
						ports []dockerxmpp.Ports
					)

					for _, port := range ap.PortsDST {
						p, _ := strconv.Atoi(port)
						ports = append(ports, dockerxmpp.Ports{Port: p})
					}

					time.Sleep(3 * time.Second)
					if elements, err = generateDeploy(dockerxmpp.Elements{
						Name:   c.Name,
						Cpus:   ap.Cpus,
						Memory: ap.Memory,
						Image:  c.Image,
						Ports:  ports,
					}); err != nil {
						config.EnvSingleton.Log.Errorf(log.TEMPLATE_ACTION, "Core", "generateDeploy", "Error to generate die container", err.Error())
						containersDie <- c
						break
					}

					for _, port := range elements.PortsContainer {
						container.Ports[port.Source] = port.Destinations
					}

					container.Address = elements.Address
					message.Subject = "Generate new container die with success"
				} else {
					message.Subject = "Container Die"
				}

				if len(masterNode) > 0 {
					if body, err = json.Marshal(container); err != nil {
						config.EnvSingleton.Log.Errorf(log.TEMPLATE_ACTION, "Core", "checkContainerDie", "Marshal Error", err.Error())
						break
					}

					message.To = masterNode[0]
					message.Body = string(body)

					if err = config.EnvSingleton.XmppConnection.Send(message); err != nil {
						config.EnvSingleton.Log.Errorf(log.TEMPLATE_ACTION, "Core", "checkContainerDie", "Error use method Send of xmpp", err.Error())
					}
				}
			}
		}
	}()
}

func initXMPP(ctx context.Context) error {
	var err error

	go config.EnvSingleton.XmppConnection.Receiver(ctx)
	config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.MESSAGE_HANDLER, Message)
	config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.PRESENCE_HANDLER, Presence)
	config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

	if err = config.EnvSingleton.XmppConnection.Roster(); err != nil {
		return err
	}

	//config.EnvSingleton.XmppConnection.DiscoItems("conference.localhost")
	//config.EnvSingleton.XmppConnection.DiscoItems("minions@conference.localhost")
	if err = config.EnvSingleton.XmppConnection.DiscoItems(config.EnvXmpp.Room); err != nil {
		return err
	}

	if err = config.EnvSingleton.XmppConnection.MucPresence(config.EnvXmpp.Room); err != nil {
		return err
	}

	return nil
}

func Message(i interface{}) {
	fmt.Println(i)
}

func Presence(i interface{}) {
	var (
		p      = i.(*xmpp.Presence)
		idx    int
		exists bool
	)

	if !reflect.DeepEqual(p.User, xmpp.MucUser{}) && strings.Contains(p.From, config.EnvXmpp.MasterUser) {
		for _, item := range p.User.Item {
			if item.Jid != EMPTY_STR {
				if p.Type == UNAVAILABLE {
					config.EnvSingleton.Log.Infof(log.TEMPLATE_PRESENCE, item.Jid, UNAVAILABLE)
					if idx, exists = utils.ExistsStringElement(item.Jid, masterNode); exists {
						masterNode = append(masterNode[idx:], masterNode[:idx+1]...)
					}
				} else {
					config.EnvSingleton.Log.Infof(log.TEMPLATE_PRESENCE, item.Jid, AVAILABLE)
					if _, exists = utils.ExistsStringElement(item.Jid, masterNode); !exists {
						masterNode = append(masterNode, item.Jid)
					}
				}
			}

			fmt.Println(masterNode)
		}
	}
}

func Iq(i interface{}) {
	var (
		v        = i.(*xmpp.ClientIQ)
		q        dockerxmpp.QueryDocker
		query    xmpp.Query
		err      error
		elements dockerxmpp.Elements
	)

	if v.Type != "error" && len(v.Query) > 0 {
		if err = xml.Unmarshal(v.Query, &query); err != nil {
			config.EnvSingleton.Log.Errorf(log.TEMPLATE_ACTION, "Core", "Iq", "Unmarshal IQ", string(v.Query))
			return
		}

		switch query.XMLName.Space {
		case JABBER_IQ_DOCKER:
			config.EnvSingleton.Log.Debugf(log.TEMPLATE_ACTION, "Core", "Iq", "Receive Docker IQ", string(v.Query))
			if err = xml.Unmarshal(v.Query, &q); err != nil {
				config.EnvSingleton.Log.Errorf(log.TEMPLATE_ACTION, "Core", "Iq", "Unmarshal IQ Docker", err.Error())
				break
			}

			var iq = dockerxmpp.IQ{
				From: v.To,
				To:   v.From,
				Type: "result",
				ID:   v.ID,
			}

			switch q.Action {
			case EMPTY_STR:
				err = errors.New("Action is not informed")
			case dockerxmpp.GENERATE_IMAGE:
				err = generateImage(q.Elements)
			case dockerxmpp.LOAD_IMAGE:
				elements, err = loadImage(q.Elements)
			case dockerxmpp.EXISTS_IMAGE:
				elements, err = existsImage(q.Elements)
			case dockerxmpp.MASTER_DEPLOY:
				var (
					ed      = make(chan EventsDocker, 1)
					ports   dockerxmpp.Elements
					address dockerxmpp.Elements
				)

				eventsDocker[q.Elements.Name] = ed
				if elements, err = masterDeploy(q.Elements); err != nil {
					delete(eventsDocker, q.Elements.Name)
					break
				}

				go func() {
					time.Sleep(5 * time.Second)
					if _, ok := eventsDocker[q.Elements.Name]; ok {
						eventsDocker[q.Elements.Name] <- EventsDocker{}
					}
				}()

				select {
				case r := <-ed:
					err = r.Error
				}
				delete(eventsDocker, q.Elements.Name)

				if address, err = addressContainer(q.Elements); err != nil {
					break
				}

				if ports, err = portsContainer(q.Elements); err != nil {
					break
				}

				elements.PortsContainer = ports.PortsContainer
				elements.Address = address.Address
				elements.Minion = config.EnvConfig.Hostname

				if err != nil {
					var c []docker.Containers
					if c, err = docker.ListAllContainers(q.Elements.Name); err == nil {
						containers = append(containers, c...)
					} else {
						containers = append(containers, docker.Containers{
							ID:   elements.ID,
							Name: q.Elements.Name,
						})
					}
				}
			case dockerxmpp.APPEND_DEPLOY:
				elements, err = appendDeploy(q.Elements)
			case dockerxmpp.NAME_CONTAINERS:
				elements, err = nameContainers()
			case dockerxmpp.TOTAL_CONTAINERS:
				elements, err = totalContainers()
			case dockerxmpp.OPERATION_CONTAINERS:
				elements, err = operationContainers(q.Elements)
			case dockerxmpp.REMOVE_CONTAINER:
				containersRemove = append(containersRemove, q.Elements.Name)
				err = removeContainer(q.Elements)

				if err == nil {
					for idx, c := range containers {
						if c.Name == q.Elements.Name {
							if len(containers)-1 == idx {
								containers = containers[:idx]
							} else {
								containers = append(containers[:idx], containers[idx+1:]...)
							}
						}
					}
				}

				var idx int
				idx, _ = utils.ExistsStringElement(q.Elements.Name, containersRemove)

				if len(containersRemove)-1 == idx {
					containersRemove = containersRemove[:idx]
				} else {
					containersRemove = append(containersRemove[:idx], containersRemove[idx+1:]...)
				}
			default:
				err = errors.New("Action is not exists")
			}

			if err == nil {
				iq.Type = "result"
				iq.Query = dockerxmpp.QueryDocker{
					Action:   q.Action,
					Elements: elements,
				}
			} else {
				iq.Type = "error"
				iq.Error = &dockerxmpp.IQError{
					Type: "cancel",
					Text: err.Error(),
				}
			}

			fmt.Println(iq.To)
			fmt.Println("ERROR: ", iq.Error)
			if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
				config.EnvSingleton.Log.Errorf(log.TEMPLATE_ACTION, "Core", "Iq", "Send message XMPP", err.Error())
			}
		}
	}
}

func generateDeploy(elements dockerxmpp.Elements) (dockerxmpp.Elements, error) {
	var (
		ed      = make(chan EventsDocker, 1)
		ports   dockerxmpp.Elements
		address dockerxmpp.Elements
		err     error
		el      = dockerxmpp.Elements{}
	)

	eventsDocker[elements.Name] = ed
	defer func() {
		delete(eventsDocker, elements.Name)
	}()

	if el, err = masterDeploy(elements); err != nil {
		return dockerxmpp.Elements{}, err
	}

	go func() {
		time.Sleep(5 * time.Second)
		if _, ok := eventsDocker[elements.Name]; ok {
			eventsDocker[elements.Name] <- EventsDocker{}
		}
	}()

	select {
	case r := <-ed:
		if r.Error != nil {
			return dockerxmpp.Elements{}, r.Error
		}
	}

	if address, err = addressContainer(elements); err != nil {
		return dockerxmpp.Elements{}, err
	}

	if ports, err = portsContainer(elements); err != nil {
		return dockerxmpp.Elements{}, err
	}

	var c []docker.Containers
	if c, err = docker.ListAllContainers(elements.Name); err == nil {
		containers = append(containers, c...)
	} else {
		containers = append(containers, docker.Containers{
			ID:   el.ID,
			Name: elements.Name,
		})
	}

	return dockerxmpp.Elements{
		ID:             el.ID,
		Name:           elements.Name,
		PortsContainer: ports.PortsContainer,
		Address:        address.Address,
		Minion:         config.EnvConfig.Hostname,
	}, nil
}
