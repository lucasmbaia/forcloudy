package config

import (
	"github.com/lucasmbaia/go-xmpp"
)

var (
	EnvXmpp      Xmpp
	EnvSingleton Singleton
)

type Config struct {
	UserMasterNode string `json:",omitempty"`
}

type Xmpp struct {
	Host                  string `json:"host,omitempty"`
	Port                  string `json:"port,omitempty"`
	MechanismAuthenticate string `json:"mechanismAuthenticate,omitempty"`
	User                  string `json:"user,omitempty"`
	Password              string `json:"password,omitempty"`
	Room                  string `json:"room,omitempty"`
}

type Singleton struct {
	XmppConnection *xmpp.Client
}

func LoadConfig() {
	loadXMPP()
}
