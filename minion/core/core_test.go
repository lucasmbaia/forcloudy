package core

import (
  "context"
  "forcloudy/minion/config"
  "github.com/lucasmbaia/go-xmpp"
  dockerxmpp "github.com/lucasmbaia/go-xmpp/docker"
  "testing"
  "encoding/xml"
)

func loadConfig() {
  config.EnvXmpp = config.Xmpp{
    Host: "172.16.95.179",
    Port: "5222",
    MechanismAuthenticate: xmpp.PLAIN,
    User:     "minion-1@localhost",
    Password: "totvs@123",
    Room:     "minions@conference.localhost",
  }

  config.LoadConfig()
}

func TestRunCore(t *testing.T) {
  loadConfig()
  Run(context.Background())
}

func TestIq(t *testing.T) {
  loadConfig()

  var iq = &xmpp.ClientIQ{
    Query:  []byte("<query xmlns=\"jabber:iq:docker\" action=\"master-deploy\"><name>teste</name><customer>lucas</customer></query> <nil>"),
  }

  Iq(iq)
}

func TestIqMasterDeploy(t *testing.T) {
  loadConfig()
  watchEvents(context.Background())

  var (
    body  []byte
    err	  error
  )

  var elements = dockerxmpp.Elements {
    Name:	  "bematech",
    Cpus:	  "0.1",
    Memory:	  "15MB",
    Image:	  "alpine",
    Ports:	  []dockerxmpp.Ports{
      {Port: 80},
    },
    CreateImage:  true,
  }

  var dx = dockerxmpp.QueryDocker{Action: dockerxmpp.MASTER_DEPLOY, Elements: elements}

  if body, err = xml.Marshal(dx); err != nil {
    t.Fatal(err)
  }

  var iq = &xmpp.ClientIQ{
    Query:  body,
  }

  Iq(iq)
}
