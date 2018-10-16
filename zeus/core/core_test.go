package core

import (
  "context"
  "forcloudy/zeus/config"
  "github.com/lucasmbaia/go-xmpp"
  "testing"
)

func loadConfig() {
  config.EnvXmpp = config.Xmpp{
    Host: "172.16.95.179",
    Port: "5222",
    MechanismAuthenticate: xmpp.PLAIN,
    User:     "zeus@localhost",
    Password: "totvs@123",
    Room:     "minions@conference.localhost",
  }

  config.EnvConfig = config.Config{
    EtcdEndpoints:  []string{"http://127.0.0.1:2379"},
    EtcdTimeout:	  5,
  }

  config.LoadConfig()
}

func TestRunCore(t *testing.T) {
  loadConfig()
  Run(context.Background())
}
