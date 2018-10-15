package config

import (
	"github.com/lucasmbaia/go-xmpp"
	"log"
)

func loadXMPP() {
	var err error
	log.Println("loadXMPP")

	if EnvSingleton.XmppConnection, err = xmpp.NewClient(xmpp.Options{
		Host:      EnvXmpp.Host,
		Port:      EnvXmpp.Port,
		Mechanism: EnvXmpp.MechanismAuthenticate,
		User:      EnvXmpp.User,
		Password:  EnvXmpp.Password,
	}); err != nil {
		log.Fatal(err)
	}
}
