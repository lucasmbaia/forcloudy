package config

import (
	"github.com/lucasmbaia/forcloudy/minion/log"
	"github.com/lucasmbaia/go-xmpp"
	_log "log"
)

func loadXMPP() {
	var err error

	EnvSingleton.Log.Infofc(log.TEMPLATE_LOAD, "Config", "loadXMPP", EnvXmpp)

	if EnvSingleton.XmppConnection, err = xmpp.NewClient(xmpp.Options{
		Host:      EnvXmpp.Host,
		Port:      EnvXmpp.Port,
		Mechanism: EnvXmpp.MechanismAuthenticate,
		User:      EnvXmpp.User,
		Password:  EnvXmpp.Password,
	}); err != nil {
		_log.Fatal(err)
	}
}
