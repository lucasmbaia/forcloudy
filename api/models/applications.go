package models

import (
  "fmt"
  "github.com/lucasmbaia/forcloudy/api/datamodels"
  "github.com/lucasmbaia/forcloudy/api/repository"
)

type Applications struct {
  repository repository.Repositorier
}

func NewApplications(session repository.Repositorier) *Applications {
  return &Applications{repository: session}
}

func (a *Applications) Get(filters interface{}) (i interface{}, err error) {
  return i, err
}

func (a *Applications) Post(values interface{}) error {
  var (
    applications = values.(*datamodels.ApplicationsFields)
  )

  fmt.Println(applications)
  return nil
}
