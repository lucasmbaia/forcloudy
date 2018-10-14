package els

import (
  "github.com/olivere/elastic"
  "context"
)

type Client struct {
  coon	*elastic.Client
  ctx	context.Context
}

func NewClient(ctx context.Context, url string) (*Client, error) {
  var (
    client  = &Client{}
    err	    error
  )

  client.ctx = ctx
  client.coon, err = elastic.NewClient(elastic.SetURL(url))

  return client, err
}

func (c *Client) Post(index, topic string, body interface{}) error {
  var (
    err	    error
    exists  bool
  )

  if exists, err = c.coon.IndexExists(index).Do(c.ctx); err != nil {
    return err
  }

  if !exists {
    if _, err = c.coon.CreateIndex(index).Do(c.ctx); err != nil {
      return err
    }
  }

  if _, err = c.coon.Index().Index(index).Type(topic).BodyJson(body).Do(c.ctx); err != nil {
    return err
  }

  return nil
}
