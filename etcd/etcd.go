package etcd

import (
	"context"
	"encoding/json"
	"errors"
	_etcd "github.com/etcd-io/etcd/client"
	"reflect"
	"time"
)

type Config struct {
	Username  string
	Password  string
	Endpoints []string
	Timeout   int32
}

type Client struct {
	cli     _etcd.KeysAPI
	timeout int32
	ctx     context.Context
}

type Broker interface {
	Set(key string, value interface{}) error
	Get(key string, value interface{}) error
}

func NewClient(ctx context.Context, config Config) (Client, error) {
	var (
		client _etcd.Client
		err    error
	)

	if client, err = _etcd.New(_etcd.Config{
		Username:                config.Username,
		Password:                config.Password,
		Endpoints:               config.Endpoints,
		Transport:               _etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Duration(config.Timeout) * time.Second,
	}); err != nil {
		return Client{}, err
	}

	return Client{
		cli:     _etcd.NewKeysAPI(client),
		timeout: config.Timeout,
		ctx:     ctx,
	}, nil
}

func (c Client) Set(key string, value interface{}) error {
	var (
		body []byte
		err  error
	)

	if body, err = json.Marshal(value); err != nil {
		return err
	}

	if _, err = c.cli.Set(c.ctx, key, string(body), nil); err != nil {
		return err
	}

	return nil
}

func (c Client) Get(key string, value interface{}) error {
	var (
		err      error
		response *_etcd.Response
	)

	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return errors.New("Expected a pointer to a variable")
	}

	if response, err = c.cli.Get(c.ctx, key, nil); err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(response.Node.Value), value); err != nil {
		return err
	}

	return nil
}
