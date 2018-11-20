package haproxy

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
	if conf, err := httpAndHttps(infos{
		ApplicationName:   "httpAndHttps",
		ContainerName:     "lucas_app-httpAndHttps-1",
		PortSource:        "80",
		PortsDestionation: []string{"32987"},
		AddressContainer:  "127.0.0.1",
		Dns:               "httpAndHttps.local",
		Minion:            "minion-1",
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

func Test_TcpAndUdp(t *testing.T) {
	if conf, err := tcpAndUdp(infos{
		Customer:          "lucas",
		ApplicationName:   "tcpAndUdp",
		ContainerName:     "lucas_app-tcpAndUdp-1",
		PortSource:        "5466",
		PortsDestionation: []string{"32987"},
		AddressContainer:  "127.0.0.1",
		Dns:               "tcpAndUdp.local",
		Protocol:          "tcp",
		Minion:            "minion-1",
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

func Test_GenerateConf(t *testing.T) {
	if err := GenerateConf(Haproxy{
		Customer:         "lucas",
		ApplicationName:  "teste",
		ContainerName:    "teste-1",
		PortsContainer:   map[string][]string{"5233": []string{"34999"}},
		Protocol:         map[string]string{"5233": "tcp"},
		AddressContainer: "127.0.0.1",
		Dns:              "teste.local",
		Minion:           "minion-1",
	}); err != nil {
		t.Fatal(err)
	}
}

func Test_RemoveContainer(t *testing.T) {
	if err := RemoveContainer(Haproxy{
		Customer:        "lucas",
		ApplicationName: "cuzao",
		ContainerName:   "lucas_app-cuzao-4",
		PortsContainer:  map[string][]string{"80": []string{"32797"}},
	}); err != nil {
		t.Fatal(err)
	}
}
