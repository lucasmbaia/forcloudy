package config

import (
	"github.com/lucasmbaia/forcloudy/logging"
	"github.com/lucasmbaia/go-xmpp"
	_log "log"
	"os"
)

var (
	EnvXmpp      Xmpp
	EnvSingleton Singleton
	EnvConfig    Config
)

type Config struct {
	UserMasterNode string `json:",omitempty"`
	Hostname       string
}

type Xmpp struct {
	Host                  string `json:"host,omitempty"`
	Port                  string `json:"port,omitempty"`
	MechanismAuthenticate string `json:"mechanismAuthenticate,omitempty"`
	User                  string `json:"user,omitempty"`
	Password              string `json:"password,omitempty"`
	Room                  string `json:"room,omitempty"`
	MasterUser            string `zeus:"masterUser,omitempty"`
}

type Singleton struct {
	XmppConnection *xmpp.Client
	Log            *logging.Logger
}

func LoadConfig() {
	var err error

	LoadLog(logging.INFO)
	loadXMPP()

	if EnvConfig.Hostname, err = os.Hostname(); err != nil {
		_log.Panic(err)
	}
}
