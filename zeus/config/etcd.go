package config

import (
  "forcloudy/etcd"
  "context"
  "log"
)

func loadETCD() {
  var err error
  log.Println("loadETCD")

  if EnvSingleton.EtcdConnection, err = etcd.NewClient(context.Background(), etcd.Config{
    Endpoints:	EnvConfig.EtcdEndpoints,
    Timeout:	EnvConfig.EtcdTimeout,
  }); err != nil {
    log.Fatal(err)
  }
}
