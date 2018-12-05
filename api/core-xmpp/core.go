package core

import (
  "github.com/lucasmbaia/forcloudy/api/config"
  "github.com/lucasmbaia/forcloudy/api/utils"
  "github.com/lucasmbaia/forcloudy/dtos"
  "github.com/lucasmbaia/go-xmpp"
  "github.com/lucasmbaia/go-xmpp/docker"

  "context"
  "encoding/json"
  "encoding/xml"
  "errors"
  "fmt"
  "reflect"
  "strconv"
  "strings"
  "sync"
)

const (
  EMPTY_STR        = ""
  JABBER_IQ_DOCKER = "jabber:iq:docker"
  UNAVAILABLE      = "unavailable"
)

var (
  mutex	    = &sync.RWMutex{}
  minions	    map[string]Minions
  masters	    []string
  responseIQ  map[string]chan ResponseIQ
)

type Minions struct {
  Containers []string
}

type ResponseIQ struct {
  Error    error
  Elements docker.Elements
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

  //go config.EnvSingleton.XmppConnection.Receiver(ctx)
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.MESSAGE_HANDLER, Message)
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.PRESENCE_HANDLER, Presence)
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

  if err = config.EnvSingleton.XmppConnection.Roster(); err != nil {
    return err
  }

  if err = config.EnvSingleton.XmppConnection.DiscoItems(config.EnvXmpp.Room); err != nil {
    return err
  }

  if err = config.EnvSingleton.XmppConnection.MucPresence(config.EnvXmpp.Room); err != nil {
    return err
  }

  if err = config.EnvSingleton.XmppConnection.MucPresence(config.EnvXmpp.MasterRoom); err != nil {
    return err
  }

  checkDown()
  return nil
}

func checkDown() {
  go func() {
    for {
      select {
      case _ = <-config.EnvXmpp.SystemShutdown:
	config.EnvSingleton.XmppConnection.Roster()
	config.EnvSingleton.XmppConnection.DiscoItems(config.EnvXmpp.Room)
	config.EnvSingleton.XmppConnection.MucPresence(config.EnvXmpp.Room)
	config.EnvSingleton.XmppConnection.MucPresence(config.EnvXmpp.MasterRoom)
      }
    }
  }()
}

func Message(i interface{}) {
  var m = i.(*xmpp.Message)

  switch m.Subject {
  case "Generate new container die with success":
    var (
      container dtos.Container
      err       error
      h         utils.Haproxy
      ap        dtos.ApplicationEtcd
    )

    if err = json.Unmarshal([]byte(m.Body), container); err != nil {
      fmt.Println(err)
      return
    }

    if config.EnvSingleton.EtcdConnection.Get(fmt.Sprintf("%s%s", utils.KEY_ETCD, container.Customer), &ap); err != nil {
      fmt.Println(err)
      return
    }

    h = utils.Haproxy{
      Customer:         container.Customer,
      ApplicationName:  container.Application,
      ContainerName:    container.Name,
      PortsContainer:   container.Ports,
      Protocol:         ap.Protocol,
      AddressContainer: container.Address,
      Dns:              ap.Dns,
      Minion:           container.Minion,
    }

    if err = utils.RemoveContainer(h); err != nil {
      fmt.Println(err)
      return
    }

    if err = utils.GenerateConf(h); err != nil {
      fmt.Println(err)
      return
    }
  case "Container Die":
    var (
      container dtos.Container
      err       error
      ap        dtos.ApplicationEtcd
      response  = make(chan Container)
      ports     []Ports
      iterator  int
    )

    if err = json.Unmarshal([]byte(m.Body), container); err != nil {
      fmt.Println(err)
      return
    }

    if config.EnvSingleton.EtcdConnection.Get(fmt.Sprintf("%s%s", utils.KEY_ETCD, container.Customer), &ap); err != nil {
      fmt.Println(err)
      return
    }

    if err = utils.RemoveContainer(utils.Haproxy{
      Customer:        container.Customer,
      ApplicationName: container.Application,
      ContainerName:   container.Name,
      PortsContainer:  container.Ports,
    }); err != nil {
      fmt.Println(err)
      return
    }

    for port, protocol := range ap.Protocol {
      p, _ := strconv.Atoi(port)
      ports = append(ports, Ports{Port: p, Protocol: protocol})
    }

    var aux = strings.Split(container.Name, "-")
    iterator, _ = strconv.Atoi(aux[len(aux)-1])

    go func() {
      select {
      case resp := <-response:
	if resp.Error != nil {
	  fmt.Println(resp.Error)
	  break
	}

	var portsContainer = make(map[string][]string)

	for _, port := range resp.PortsContainer {
	  portsContainer[port.Source] = port.Destinations
	}

	if err = utils.GenerateConf(utils.Haproxy{
	  Customer:         container.Customer,
	  ApplicationName:  container.Application,
	  ContainerName:    container.Name,
	  PortsContainer:   portsContainer,
	  Protocol:         ap.Protocol,
	  AddressContainer: resp.Address,
	  Dns:              ap.Dns,
	  Minion:           resp.Minion,
	}); err != nil {
	  fmt.Println(err)
	  break
	}
      }
    }()

    if err = DeployApplication(Deploy{
      Customer:        container.Customer,
      ApplicationName: container.Application,
      Cpus:            ap.Cpus,
      Memory:          ap.Memory,
      ImageVersion:    strings.Split(container.Image, ":")[1],
      TotalContainers: 1,
      Ports:           ports,
    }, iterator, false, response); err != nil {
      fmt.Println(err)
      return
    }
  }
  fmt.Println(i)
}

func Presence(i interface{}) {
  var p = i.(*xmpp.Presence)

  fmt.Println(p)
  if !reflect.DeepEqual(p.User, xmpp.MucUser{}) && strings.Contains(p.From, config.EnvXmpp.MasterRoom) {
    if strings.Replace(p.From, fmt.Sprintf("%s/", config.EnvXmpp.MasterRoom), "", -1) != config.EnvSingleton.XmppConnection.User {
      var (
	idx	  int
	exists  bool
      )

      idx, exists = utils.ExistsStringElement(p.From, masters)

      if p.Type == UNAVAILABLE {
	if exists {
	  if len(masters) -1 == idx {
	    masters = masters[:idx]
	  } else {
	    masters = append(masters[:idx], masters[idx + 1:]...)
	  }
	}
      } else {
	if !exists {
	  masters = append(masters, p.From)
	}
      }
    }

    fmt.Println(masters)
  }

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
	    iq         docker.IQ
	    err        error
	    response   = make(chan ResponseIQ, 1)
	    containers []string
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

	      if err = config.EnvSingleton.XmppConnection.SendPresenceMuc(item.Jid); err != nil {
		return
	      }
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
    iq    = i.(*xmpp.ClientIQ)
    query xmpp.Query
    err   error
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
