package core

import (
	"context"
	"forcloudy/minion/config"
	"github.com/lucasmbaia/go-xmpp"
	"testing"
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