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

	config.LoadConfig()
}

func TestRunCore(t *testing.T) {
	loadConfig()
	Run(context.Background())
}
