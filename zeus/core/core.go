package core

import (
  "context"
  "fmt"
  "encoding/xml"
  "forcloudy/zeus/config"
  "github.com/lucasmbaia/go-xmpp"
  "github.com/lucasmbaia/go-xmpp/docker"
  "log"
  "reflect"
  "strings"
  "sync"
)

const (
  EMPTY_STR        = ""
  JABBER_IQ_DOCKER = "jabber:iq:docker"
)

var (
  mutex	      = &sync.RWMutex{}
  minions     map[string]Minions
  responseIQ  map[string]chan ResponseIQ
)

type Minions struct {
  Containers []string
}

type ResponseIQ struct {
  Error	    error
  Elements  docker.Elements
}

func init() {
  minions = make(map[string]Minions)
  responseIQ = make(map[string]chan ResponseIQ)
}

func Run(ctx context.Context) error {
  var (
    err error
  )

  if err = initXMPP(ctx); err != nil {
    return err
  }

  select {
  case _ = <-ctx.Done():
    return nil
  }
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
  var (
    msg	= i.(*xmpp.Message)
  )

  if msg.Body != EMPTY_STR {
    fmt.Println(msg)
  }
}

func Presence(i interface{}) {
  var v = i.(*xmpp.Presence)

  if !reflect.DeepEqual(v.User, xmpp.MucUser{}) && !strings.Contains(v.From, config.EnvSingleton.XmppConnection.User) && strings.Contains(v.From, config.EnvXmpp.Room) {
    for _, item := range v.User.Item {
      if v.Type == "unavailable" {
	mutex.Lock()
	if _, ok := minions[item.Jid]; ok {
	  delete(minions, item.Jid)
	}
	mutex.Unlock()
      } else {
	mutex.Lock()
	if _, ok := minions[item.Jid]; !ok {
	  var (
	    iq		docker.IQ
	    err		error
	    response	= make(chan ResponseIQ, 1)
	    containers	[]string
	  )

	  if iq, err = docker.NameContainers(config.EnvSingleton.XmppConnection.Jid, item.Jid); err != nil {
	    log.Println(err)
	    return
	  }

	  responseIQ[iq.ID] = response

	  if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
	    log.Println(err)
	    return
	  }

	  go func() {
	    select {
	    case r := <-response:
	      for _, container := range r.Elements.Containers {
		containers = append(containers, container.Name)
	      }

	      minions[item.Jid] = Minions{Containers: containers}
	    }
	  }()
	}
	mutex.Unlock()
      }
    }
  }

  fmt.Println("PRESENCE")
  fmt.Println(v)
}

func Iq(i interface{}) {
  var (
    v         = i.(*xmpp.ClientIQ)
    q         docker.QueryDocker
    query     xmpp.Query
    err       error
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

      var response = ResponseIQ{Elements: q.Elements}
      mutex.Lock()
      if _, ok := responseIQ[v.ID]; ok {
	responseIQ[v.ID] <- response
      }
      mutex.Unlock()
    }
  }

  fmt.Println(string(v.Query))
}
