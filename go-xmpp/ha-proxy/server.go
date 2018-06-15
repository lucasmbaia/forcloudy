package main

import (
  "context"
  "os"
  "os/signal"
  "syscall"
  "forcloudy/go-xmpp/ha-proxy/watch"
  "flag"
  "log"
  "fmt"
)

var (
  timeout = flag.Int("timeout", 5, "timeout of connect etcd")
  key = flag.String("key", "/haproxy/", "Key of watch etcd")
  hosts = flag.String("host", "http://172.16.95.183:2379", "Host of etcd")
)

type InfosApplication struct {
  Name  string  `json:"name,omitempty"`
  Hosts []Hosts `json:"hosts,omitempty"`
  Dns   string  `json:"dns,omitempty"`
}

type Hosts struct {
  PortSRC string   `json:"portSRC,omitempty"`
  Address []string `json:"address,omitempty"`
}

func main() {
  var (
    ctx	    context.Context
    cancel  context.CancelFunc
    err	    error
    cli	    *watch.Client
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

      fmt.Println(infos.Key, infos.Values)
    }
  }()

  go func() {
    if err = cli.Watch(*key, values); err != nil {
      log.Fatalf("Watch error: %s", err.Error())
    }
  }()

  <-ctx.Done()
}
