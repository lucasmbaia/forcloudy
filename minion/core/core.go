package core

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"forcloudy/minion/config"
	"forcloudy/minion/docker"
	"github.com/lucasmbaia/go-xmpp"
	dockerxmpp "github.com/lucasmbaia/go-xmpp/docker"
	"log"
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
)

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

	//config.EnvSingleton.XmppConnection.DiscoItems("conference.localhost")
	//config.EnvSingleton.XmppConnection.DiscoItems("minions@conference.localhost")

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
		v   = i.(*xmpp.ClientIQ)
		q   xmpp.Query
		qd  = dockerxmpp.QueryDocker{}
		err error
	)

	if v.Type != "error" {
		if err = xml.Unmarshal(v.Query, &q); err != nil {
			log.Println(err)
		}

		switch q.XMLName.Space {
		case JABBER_IQ_DOCKER:
			fmt.Println("TOMA NO CU")
			if err = xml.Unmarshal(v.Query, &qd); err != nil {
				log.Println(err)
			}

			fmt.Println("AQUI")
			fmt.Println(qd)
		}
		fmt.Println(q.XMLName.Space)
	}

	fmt.Println(string(v.Query))
}
