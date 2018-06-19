package main

import (
  "context"
  "os"
  "os/signal"
  "syscall"
  "forcloudy/go-xmpp/ha-proxy/watch"
  "forcloudy/go-xmpp/ha-proxy/template"
  "encoding/json"
  "strings"
  "flag"
  "log"
  "fmt"
)

var (
  timeout = flag.Int("timeout", 5, "timeout of connect etcd")
  key = flag.String("key", "/haproxy/", "Key of watch etcd")
  hosts = flag.String("host", "http://172.16.95.183:2379", "Host of etcd")
  path = flag.String("path", "", "Path to conf ha-proxy")
)

type InfosApplication struct {
  Name  string  `json:"name,omitempty"`
  Hosts []Hosts `json:"hosts,omitempty"`
  Dns   string  `json:"dns,omitempty"`
}

type Hosts struct {
  Protocol  string    `json:"protocol,omitempty"`
  PortSRC   string    `json:"portSRC,omitempty"`
  Address   []string  `json:"address,omitempty"`
  Whitelist string    `json:"-"`
}

func Whitelist(address []string) string {
  var addrs string

  for _, v := range address {
    addrs = fmt.Sprintf("%s%s ", addrs, strings.Split(v, ":")[0])
  }

  addrs = fmt.Sprintf("%s%s %s %s", addrs, "minion-1", "minion-2", "minion-3")
  return addrs
  //return addrs[:len(addrs) - 1]
}

func main() {
  var (
    ctx	    context.Context
    cancel  context.CancelFunc
    err	    error
    cli	    *watch.Client
    ia	    InfosApplication
    values  = make(chan watch.WatchInfos)
    sigs    = make(chan os.Signal, 1)
  )

  ctx, cancel = context.WithCancel(context.Background())

  signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

  go func() {
    <-sigs
    cancel()
  }()

  if cli, err = watch.New([]string{*hosts}, *timeout); err != nil {
    log.Fatalf("Error to connect etcd: %s", err.Error())
  }

  go func() {
    for {
      infos := <-values

      if err = json.Unmarshal([]byte(infos.Values), &ia); err != nil {
	log.Printf("Error unmarshal: %s", err.Error())
	continue
      }

      fmt.Println(infos.Key, infos.Values)

      for key, host := range ia.Hosts {
	ia.Hosts[key].Whitelist = Whitelist(host.Address)
      }

      ia.Name = strings.Replace(infos.Key, *key, "", 1)
      if err = template.ConfGenerate(*path, ia.Name, template.MINION, ia); err != nil {
	log.Printf("Error to generate conf: %s", err.Error())
	continue
      }
    }
  }()

  go func() {
    if err = cli.Watch(*key, values); err != nil {
      log.Fatalf("Watch error: %s", err.Error())
    }
  }()

  <-ctx.Done()
}
