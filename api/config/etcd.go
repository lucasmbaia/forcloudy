package config

import (
	"context"
	"github.com/lucasmbaia/forcloudy/etcd"
)

func LoadETCD() {
	var err error

	if EnvSingleton.EtcdConnection, err = etcd.NewClient(context.Background(), etcd.Config{
		Endpoints: EnvConfig.EtcdEndpoints,
		Timeout:   EnvConfig.EtcdTimeout,
	}); err != nil {
		panic(err)
	}
}
