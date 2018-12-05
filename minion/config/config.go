package config

import (
	"github.com/lucasmbaia/forcloudy/etcd"
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
	EtcdUsername   string
	EtcdPassword   string
	EtcdEndpoints  []string
	EtcdTimeout    int32
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
	EtcdConnection etcd.Client
}

func LoadConfig() {
	var err error

	LoadLog(logging.INFO)
	loadXMPP()
	//LoadETCD()

	if EnvConfig.Hostname, err = os.Hostname(); err != nil {
		_log.Panic(err)
	}
}
