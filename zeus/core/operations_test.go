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
  loadConfig()
  go config.EnvSingleton.XmppConnection.Receiver(context.Background())
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

  minions["minion-1@localhost/1085312404117464703924194"] = Minions{}
  minions["minion-2@localhost/1600762215542285951024258"] = Minions{}

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
    Build:	      "hello_world",
  }, 1, false); err != nil {
    t.Fatal(err)
  }
}

func TestExistsImage(t *testing.T) {
  loadConfig()
  go config.EnvSingleton.XmppConnection.Receiver(context.Background())
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

  if _, err := existsImage("lucas/bematech:v1", "minion-1@localhost/642059527718508167717858"); err != nil {
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
  }, "minion-1@localhost/1811322184973374488622978", "lucas_app-bematech", "alpine", true); err != nil {
    t.Fatal(err)
  }

  fmt.Println(id)
}

func TestGenerateImage(t *testing.T) {
  loadConfig()
  go config.EnvSingleton.XmppConnection.Receiver(context.Background())
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

  if err := generateImage("lucas_app-bematech", "v1", "hello_world", "minion-1@localhost/1099900166375465703618306"); err != nil {
    t.Fatal(err)
  }
}

func TestRemoveContainer(t *testing.T) {
  loadConfig()
  go config.EnvSingleton.XmppConnection.Receiver(context.Background())
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

  if err := removeContainer("bematech", "minion-1@localhost/859336255431168411417154"); err != nil {
    t.Fatal(err)
  }
}

func TestTotalContainersMinion(t *testing.T) {
  loadConfig()
  go config.EnvSingleton.XmppConnection.Receiver(context.Background())
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

  var (
    total int
    err	  error
  )

  if total, err = totalContainersMinion("minion-1@localhost/524203912100248867123170"); err != nil {
    t.Fatal(err)
  }

  fmt.Println(total)
}

func TestContainersPerMinion(t *testing.T) {
  loadConfig()
  go config.EnvSingleton.XmppConnection.Receiver(context.Background())
  config.EnvSingleton.XmppConnection.RegisterHandler(xmpp.IQ_HANDLER, Iq)

  //minions["minion-1@localhost/1142267004332576407523874"] = Minions{}
  //minions["minion-2@localhost/109782955070123812623938"] = Minions{}

  var (
    total map[string]int
    err	  error
  )

  if total, err = containersPerMinion(3); err != nil {
    t.Fatal(err)
  }

  fmt.Println(total)
}
