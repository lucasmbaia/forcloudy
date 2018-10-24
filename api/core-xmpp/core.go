package core

import (
  "github.com/lucasmbaia/go-xmpp"
  "github.com/lucasmbaia/go-xmpp/docker"
  "github.com/lucasmbaia/forcloudy/api/config"

  "encoding/xml"
  "reflect"
  "strings"
  "context"
  "errors"
  "sync"
  "fmt"
)

const (
  EMPTY_STR	    = ""
  JABBER_IQ_DOCKER  = "jabber:iq:docker"
  UNAVAILABLE	    = "unavailable"
)

var (
  mutex	      = &sync.RWMutex{}
  minions     map[string]Minions
  responseIQ  map[string]chan ResponseIQ
)

type Minions struct {
  Containers  []string
}

type ResponseIQ struct {
  Error	    error
  Elements  docker.Elements
}

func init() {
  minions     = make(map[string]Minions)
  responseIQ  = make(map[string]chan ResponseIQ)
}

func Run(ctx context.Context) error {
  var (
    err	error
  )

  if err = initXMPP(ctx, config.EnvXmpp.Room); err != nil {
    return err
  }

  select {
  case _ = <-ctx.Done():
    return nil
  }
}

func initXMPP(ctx context.Context, chat string) error {
  var err error

  go config.EnvSingleton.XmppConnection.Receiver(ctx)
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.MESSAGE_HANDLER, Message)
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.PRESENCE_HANDLER, Presence)
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

  if err = config.EnvSingleton.XmppConnection.Roster(); err != nil {
    return err
  }

  if err = config.EnvSingleton.XmppConnection.DiscoItems(chat); err != nil {
    return err
  }

  if err = config.EnvSingleton.XmppConnection.MucPresence(chat); err != nil {
    return err
  }

  return nil
}

func Message(i interface{}) {
  fmt.Println(i)
}

func Presence(i interface{}) {
  var p = i.(*xmpp.Presence)

  if !reflect.DeepEqual(p.User, xmpp.MucUser{}) && !strings.Contains(p.From, config.EnvSingleton.XmppConnection.User) && strings.Contains(p.From, config.EnvXmpp.Room) {
    for _, item := range p.User.Item {
      if p.Type == UNAVAILABLE {
	mutex.Lock()
	if _, ok := minions[item.Jid]; ok {
	  delete(minions, item.Jid)
	}
	mutex.Unlock()
      } else {
	mutex.Lock()
	if _, ok := minions[item.Jid]; !ok {
	  var (
	    iq          docker.IQ
	    err         error
	    response    = make(chan ResponseIQ, 1)
	    containers  []string
	  )

	  if iq, err = docker.NameContainers(config.EnvSingleton.XmppConnection.Jid, item.Jid); err != nil {
	    return
	  }

	  responseIQ[iq.ID] = response

	  if err = config.EnvSingleton.XmppConnection.Send(iq); err != nil {
	    return
	  }

	  go func() {
	    select {
	    case r := <-response:
	      for _, container := range r.Elements.Containers {
		containers = append(containers, container.Name)
	      }

	      minions[item.Jid] = Minions{Containers: containers}
	      fmt.Println("MINIONS: ", minions)
	    }
	  }()
	}
	mutex.Unlock()
      }
    }
  }
}

func Iq(i interface{}) {
  var (
    iq     = i.(*xmpp.ClientIQ)
    query xmpp.Query
    err	  error
  )

  if iq.Type != "error" {
    if err = xml.Unmarshal(iq.Query, &query); err != nil {
      mutex.Lock()
      if _, ok := responseIQ[iq.ID]; ok {
	responseIQ[iq.ID] <- ResponseIQ{Error: err}
      }
      mutex.Unlock()
      return
    }

    switch query.XMLName.Space {
    case JABBER_IQ_DOCKER:
      var q docker.QueryDocker
      if err = xml.Unmarshal(iq.Query, &q); err != nil {
	break
      }

      var response = ResponseIQ{Elements: q.Elements}
      mutex.Lock()
      if _, ok := responseIQ[iq.ID]; ok {
	responseIQ[iq.ID] <- response
      }
      mutex.Unlock()
    }
  } else {
    mutex.Lock()
    if _, ok := responseIQ[iq.ID]; ok {
      responseIQ[iq.ID] <- ResponseIQ{Error: errors.New(iq.Error.Text)}
    }
    mutex.Unlock()
  }
}
