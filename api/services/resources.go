package services

import (
  "encoding/json"
  "fmt"
  "github.com/lucasmbaia/forcloudy/api/models"
  "github.com/lucasmbaia/forcloudy/api/repository"
)

type ResourceService interface {
  Print()
  Set(params map[string]interface{}) error
  GetFields() interface{}
  Post() error
  Get() (interface{}, error)
}

type resourceService struct {
  fields     interface{}
  model      models.Models
  repository repository.Repositorier
}

func (r *resourceService) Print() {
  fmt.Println("MODEL")
  fmt.Println(r.fields)
}

func (r *resourceService) GetFields() interface{} {
  return r.fields
}

func (r *resourceService) Set(params map[string]interface{}) error {
  var (
    body []byte
    err  error
  )

  if body, err = json.Marshal(params); err != nil {
    return err
  }

  if err = json.Unmarshal(body, r.fields); err != nil {
    return err
  }

  return nil
}

func (r resourceService) Get() (interface{}, error) {
  return r.model.Get(r.fields)
}

func (r *resourceService) Post() error {
  return r.model.Post(r.fields)
}
