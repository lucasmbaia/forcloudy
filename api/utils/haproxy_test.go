package utils

import (
	"encoding/json"
	"fmt"
	"github.com/lucasmbaia/forcloudy/api/config"
	"testing"
)

func init() {
	config.EnvConfig = config.Config{
		EtcdEndpoints: []string{"http://192.168.204.128:2379"},
		EtcdTimeout:   10,
		Hostname:      "minion-1",
	}

	config.LoadETCD()
}

func Test_HttpAndHttps(t *testing.T) {
	if conf, err := httpAndHttps(httpHttps{
		ApplicationName:   "httpAndHttps",
		ContainerName:     "lucas_app-httpAndHttps-1",
		PortSource:        "80",
		PortsDestionation: []string{"32987"},
		AddressContainer:  "127.0.0.1",
		Dns:               "httpAndHttps.local",
	}); err != nil {
		t.Fatal(err)
	} else {
		if body, err := json.Marshal(conf); err != nil {
			t.Fatal(err)
		} else {
			fmt.Println(string(body))
		}
	}
}
