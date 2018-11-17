package core

import (
	"context"
	"fmt"
	"github.com/lucasmbaia/forcloudy/api/config"
	"github.com/lucasmbaia/forcloudy/api/repository/gorm"
	"github.com/lucasmbaia/go-xmpp"
	"testing"
)

func init() {
	config.EnvXmpp = config.Xmpp{
		Host: "192.168.204.129",
		Port: "5222",
		MechanismAuthenticate: "PLAIN",
		User:     "zeus@localhost",
		Password: "totvs@123",
		Room:     "minions@conference.localhost",
	}

	config.EnvDB = config.Database{
		gorm.Config{
			Username:     "forcloudy",
			Password:     "123456",
			Host:         "localhost",
			Port:         "3306",
			DBName:       "forcloudy",
			Timeout:      "10000ms",
			Debug:        true,
			ConnsMaxIdle: 10,
			ConnsMaxOpen: 10,
		},
	}

	config.LoadConfig()
}

func Test_ExistsImage(t *testing.T) {
	go config.EnvSingleton.XmppConnection.Receiver(context.Background())
	config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

	if _, err := existsImage("lucas/bematech:v1", "minion-1@localhost/26796579502127552951122"); err != nil {
		t.Fatal(err)
	}
}

func Test_CreateContainer(t *testing.T) {
	go config.EnvSingleton.XmppConnection.Receiver(context.Background())
	config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

	if container, err := createContainer(Deploy{
		Customer:        "lucas",
		ApplicationName: "bematech",
		ImageVersion:    "v1",
		Cpus:            "0.1",
		Memory:          "15MB",
		Dns:             "bematech.local",
		Ports: []Ports{
			{Port: 80, Protocol: "http"},
		},
	}, "minion-1@localhost/167364761069969351131298", "lucas_app-bematech", "alpine", true); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(container)
	}
}

func Test_GenerateImage(t *testing.T) {
	go config.EnvSingleton.XmppConnection.Receiver(context.Background())
	config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

	if err := generateImage("lucas_app-bematech", "v1", "/root/go/src/github.com/lucasmbaia/forcloudy/minion/core/", "hello_world", "minion-1@localhost/136425763789549287001282"); err != nil {
		t.Fatal(err)
	}
}

func Test_RemoveContainer(t *testing.T) {
	go config.EnvSingleton.XmppConnection.Receiver(context.Background())
	config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

	if err := removeContainer("lucas_app-bematech", "minion-1@localhost/135779595559180451361346"); err != nil {
		t.Fatal(err)
	}
}

func Test_LoadImage(t *testing.T) {
	go config.EnvSingleton.XmppConnection.Receiver(context.Background())
	config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

	if err := loadImage("/images/", "lucas_app-bematech.tar.gz", "minion-1@localhost/75777163010889353931410"); err != nil {
		t.Fatal(err)
	}
}

func Test_DeployApplication(t *testing.T) {
	go config.EnvSingleton.XmppConnection.Receiver(context.Background())
	config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

	minions["minion-1@localhost/20269030062242505682194"] = Minions{}
	minions["minion-2@localhost/13904112234493430912162"] = Minions{}

	if err := DeployApplication(Deploy{
		Customer:        "lucas",
		ApplicationName: "bematech",
		ImageVersion:    "v1",
		Cpus:            "0.1",
		Memory:          "15MB",
		TotalContainers: 12,
		Dns:             "bematech.local",
		Ports: []Ports{
			{Port: 80, Protocol: "http"},
		},
		Path:  "/root/go/src/github.com/lucasmbaia/forcloudy/minion/core/",
		Build: "hello_world",
	}, 11, true, nil); err != nil {
		t.Fatal(err)
	}
}

func Test_ContainersPerMinion(t *testing.T) {
	go config.EnvSingleton.XmppConnection.Receiver(context.Background())
	config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

	minions["minion-1@localhost/93373821005719801381938"] = Minions{}
	minions["minion-2@localhost/144900622627952315471906"] = Minions{}

	if total, err := containersPerMinion(12); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(total)
	}
}

func Test_TotalContainersMinion(t *testing.T) {
	go config.EnvSingleton.XmppConnection.Receiver(context.Background())
	config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

	if total, err := totalContainersMinion("minion-1@localhost/143264789529165471341714"); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(total)
	}
}
