package models

import (
  "github.com/lucasmbaia/forcloudy/api/repository"
)

type Resources struct {
  repository  repository.Repositorier
}

func NewResources(session repository.Repositorier) *Resources {
  return &Resources{repository: session}
}

func (r *Resources) Post(fields interface{}) error {
  return nil
}

func (r *Resources) Get(filters interface{}) (i interface{}, err error) {
  return i, err
}

func (r *Resources) Delete(conditions interface{}) error {
  return nil
}

func (r *Resources) Put(fields, data interface{}) error {
  return nil
}
