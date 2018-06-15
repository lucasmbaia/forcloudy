package watch

import (
  "context"
  "github.com/coreos/etcd/client"
  "time"
)

type Client struct {
  client client.KeysAPI
}

type WatchInfos struct {
  Key	  string
  Values  string
}

func New(hosts []string, timeout int) (*Client, error) {
  var (
    cli  client.Client
    err  error
    resp = &Client{}
  )

  if cli, err = client.New(client.Config{
    Endpoints:               hosts,
    Transport:               client.DefaultTransport,
    HeaderTimeoutPerRequest: time.Duration(timeout) * time.Second,
  }); err != nil {
    return resp, err
  }

  resp.client = client.NewKeysAPI(cli)

  return resp, nil
}

func (c *Client) Watch(key string, values chan<- WatchInfos) error {
  var (
    watch    client.Watcher
    err      error
    response *client.Response
    ctx	     = context.Background()
  )

  for {
    watch = c.client.Watcher(key, &client.WatcherOptions{Recursive: true})

    if response, err = watch.Next(ctx); err != nil {
      return err
    }

    values <- WatchInfos{Key: response.Node.Key, Values: response.Node.Value}
  }
}
