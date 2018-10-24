package config

import (
  "github.com/lucasmbaia/go-xmpp"
  "github.com/lucasmbaia/forcloudy/etcd"
  "github.com/lucasmbaia/forcloudy/api/repository/gorm"
)

var (
  EnvSingleton Singletons
  EnvXmpp      Xmpp
  EnvDB        Database
)

type Config struct {
  EtcdUsername   string
  EtcdPassword   string
  EtcdEndpoints  []string
  EtcdTimeout    int32
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

type Database struct {
  gorm.Config
}

type Singletons struct {
  XmppConnection  *xmpp.Client
  EtcdConnection  etcd.Client
  DBConnection	  *gorm.Client
}

func LoadConfig() {
  loadDB()
  loadXMPP()
}