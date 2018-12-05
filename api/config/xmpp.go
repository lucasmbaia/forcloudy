package config

import (
  "github.com/lucasmbaia/go-xmpp"
)

func loadXMPP() {
  var err error

  if EnvSingleton.XmppConnection, err = xmpp.NewClient(xmpp.Options{
    Host:      EnvXmpp.Host,
    Port:      EnvXmpp.Port,
    Mechanism: EnvXmpp.MechanismAuthenticate,
    User:      EnvXmpp.User,
    Password:  EnvXmpp.Password,
  }, EnvXmpp.SystemShutdown); err != nil {
    panic(err)
  }
}
