package core

import (
	"context"
	"forcloudy/minion/config"
	"github.com/lucasmbaia/go-xmpp"
	"testing"
)

func loadConfig() {
	config.EnvXmpp = config.Xmpp{
		Host: "192.168.204.131",
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
