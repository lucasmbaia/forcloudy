package core

import (
  "github.com/lucasmbaia/go-xmpp"
  "forcloudy/zeus/config"
  "forcloudy/etcd"
  "testing"
  "context"
  "fmt"
)

func loadEtcd() {
  var err error
  config.EnvSingleton = config.Singleton{}

  if config.EnvSingleton.EtcdConnection, err = etcd.NewClient(context.Background(), etcd.Config{
    Endpoints:  []string{"http://127.0.0.1:2379"},
    Timeout:    5,
  }); err != nil {
    panic(err)
  }
}

func TestDeployApplication(t *testing.T) {
  loadEtcd()

  if err := DeployAppication(Deploy{
    Customer:	      "lucas",
    ApplicationName:  "bematech",
    ImageVersion:     "v1",
    Cpus:	      "0.1",
    Memory:	      "15MB",
    TotalContainers:  10,
    Dns:	      "bematech.local",
    Ports:	      []Ports{
      {Port: 80, Protocol: "http"},
    },
  }, false); err != nil {
    t.Fatal(err)
  }
}

func TestExistsImage(t *testing.T) {
  loadConfig()
  go config.EnvSingleton.XmppConnection.Receiver(context.Background())
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

  if _, err := existsImage("lucas/bematech:v1", "minion-1@localhost/1629100286397283903114594"); err != nil {
    t.Fatal(err)
  }
}

func TestDeploy(t *testing.T) {
  loadConfig()
  go config.EnvSingleton.XmppConnection.Receiver(context.Background())
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

  var (
    id	string
    err	error
  )

  if id, err = deploy(Deploy{
    Customer:	      "lucas",
    ApplicationName:  "bematech",
    ImageVersion:     "v1",
    Cpus:	      "0.1",
    Memory:	      "15MB",
    Dns:	      "bematech.local",
    Ports:	      []Ports{
      {Port: 80, Protocol: "http"},
    },
  }, "minion-1@localhost/930250711652396763816834", "bematech", "alpine", true); err != nil {
    t.Fatal(err)
  }

  fmt.Println(id)
}

func TestGenerateImage(t *testing.T) {
  loadConfig()
  go config.EnvSingleton.XmppConnection.Receiver(context.Background())
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

  if err := generateImage("bematech", "v1", "hello_world", "minion-1@localhost/1508298725114910835416962"); err != nil {
    t.Fatal(err)
  }
}
