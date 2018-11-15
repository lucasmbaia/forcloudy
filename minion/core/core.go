package core

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/lucasmbaia/forcloudy/minion/config"
	"github.com/lucasmbaia/forcloudy/minion/docker"
	"github.com/lucasmbaia/go-xmpp"
	dockerxmpp "github.com/lucasmbaia/go-xmpp/docker"
	"log"
	"time"
	//"reflect"
	//"strings"
)

const (
	EMPTY_STR        = ""
	JABBER_IQ_DOCKER = "jabber:iq:docker"
)

var (
	masterNode        []string
	containers        []docker.Containers
	containers_deploy []string
	containers_die    []string
	eventsDocker      map[string]chan EventsDocker
)

type EventsDocker struct {
	Error error
}

func init() {
	eventsDocker = make(map[string]chan EventsDocker)
}

func Run(ctx context.Context) error {
	var (
		err error
	)

	if containers, err = docker.ListAllContainers(EMPTY_STR); err != nil {
		return err
	}

	//init watch events of docker
	watchEvents(ctx)

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
		err   error
		errc  = make(chan error, 1)
		event = make(chan []byte)
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
					log.Panic(err)
				}

				if ev.Status == "die" {
					if ev.Actor.Attributes.Name != EMPTY_STR {
						if _, ok := eventsDocker[ev.Actor.Attributes.Name]; ok {
							eventsDocker[ev.Actor.Attributes.Name] <- EventsDocker{Error: errors.New("Erro to create container")}
						}
					}
				}

				log.Println(ev)
			case e := <-errc:
				log.Panic(e)
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

	config.EnvSingleton.XmppConnection.DiscoItems("conference.localhost")
	config.EnvSingleton.XmppConnection.DiscoItems("minions@conference.localhost")

	if err = config.EnvSingleton.XmppConnection.MucPresence(config.EnvXmpp.Room); err != nil {
		return err
	}

	return nil
}

func Message(i interface{}) {
	fmt.Println(i)
}

func Presence(i interface{}) {
	var v = i.(*xmpp.Presence)

	/*if !reflect.DeepEqual(v.User, xmpp.MucUser{}) && !strings.Contains(v.From, "minion-1") && strings.Contains(v.From, config.EnvXmpp.Room) {
	  }*/

	fmt.Println("PRESENCE")
	fmt.Println(v)
}

func Iq(i interface{}) {
	var (
		v        = i.(*xmpp.ClientIQ)
		q        dockerxmpp.QueryDocker
		query    xmpp.Query
		err      error
		elements dockerxmpp.Elements
	)

	if v.Type != "error" {
		if err = xml.Unmarshal(v.Query, &query); err != nil {
			log.Println(err)
		}

		switch query.XMLName.Space {
		case JABBER_IQ_DOCKER:
			if err = xml.Unmarshal(v.Query, &q); err != nil {
				break
			}

			var iq = dockerxmpp.IQ{
				From: v.To,
				To:   v.From,
				Type: "result",
				ID:   v.ID,
			}

			fmt.Println(q.Action)
			switch q.Action {
			case EMPTY_STR:
				err = errors.New("Action is not informed")
			case dockerxmpp.GENERATE_IMAGE:
				fmt.Println("GENERATE IMAGE PORRA")
				err = generateImage(q.Elements)
			case dockerxmpp.LOAD_IMAGE:
				elements, err = loadImage(q.Elements)
			case dockerxmpp.EXISTS_IMAGE:
				elements, err = existsImage(q.Elements)
			case dockerxmpp.MASTER_DEPLOY:
				fmt.Println("DEPLOY: ", q.Elements.Name)
				var (
					ed = make(chan EventsDocker, 1)
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
			case dockerxmpp.APPEND_DEPLOY:
				elements, err = appendDeploy(q.Elements)
			case dockerxmpp.NAME_CONTAINERS:
				elements, err = nameContainers()
			case dockerxmpp.TOTAL_CONTAINERS:
				elements, err = totalContainers()
			case dockerxmpp.OPERATION_CONTAINERS:
				elements, err = operationContainers(q.Elements)
			case dockerxmpp.REMOVE_CONTAINER:
				err = removeContainer(q.Elements)
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
				log.Println(err)
			}
		}
	}

	fmt.Println(string(v.Query))
}
