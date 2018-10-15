package etcd

import (
  "testing"
  "context"
  "fmt"
)

func connect() (Client, error) {
  return NewClient(context.Background(), Config{
    Endpoints:	[]string{"http://127.0.0.1:2379"},
    Timeout:	5,
  })
}

func TestNewClient(t *testing.T) {
  var err error

  if _, err = NewClient(context.Background(), Config{
    Endpoints:	[]string{"http://127.0.0.1:2379"},
    Timeout:	5,
  }); err != nil {
    t.Fatal(err)
  }
}

func TestSet(t *testing.T) {
  var s = struct {
    Name      string  `json:",omitempty"`
    LastName  string  `json:",omitempty"`
  }{"teste", "etcd"}

  var client Client
  var err error

  if client, err = connect(); err != nil {
    t.Fatal(err)
  }

  if err = client.Set("test-set", s); err != nil {
    t.Fatal(err)
  }
}

func TestGet(t *testing.T) {
  var s = struct {
    Name      string  `json:",omitempty"`
    LastName  string  `json:",omitempty"`
  }{}

  var client Client
  var err error

  if client, err = connect(); err != nil {
    t.Fatal(err)
  }

  if err = client.Get("test-set", &s); err != nil {
    t.Fatal(err)
  }

  fmt.Println(s)
}
