package watch

import (
  "context"
  "testing"
  "fmt"
)

const (
  TIMEOUT = 5
  KEY	  = "/python/"
)

var hosts = []string{"http://172.16.95.183:2379"}

func Test_NewClient(t *testing.T) {
  var (
    err	error
  )

  if _, err = New(hosts, TIMEOUT); err != nil {
    t.Fatalf("Erro to connect etcd: %s", err.Error())
  }
}

func Test_WatchEtcd(t *testing.T) {
  var (
    ctx	    context.Context
    cancel  context.CancelFunc
    err	    error
    cli	    *Client
    values  = make(chan WatchInfos)
  )

  ctx, cancel = context.WithCancel(context.Background())

  if cli, err = New(hosts, TIMEOUT); err != nil {
    t.Fatalf("Error to connect etcd: %s", err.Error())
  }

  go func() {
    for {
      infos := <-values

      fmt.Println(infos.Key, infos.Values)
    }

    cancel()
  }()

  go func() {
    if err = cli.Watch(KEY, values); err != nil {
      t.Fatalf("Watch error: %s", err.Error())
    }
  }()

  <-ctx.Done()
}
