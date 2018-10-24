package models

import (
  //"fmt"
  "github.com/lucasmbaia/forcloudy/api/datamodels"
  "github.com/lucasmbaia/forcloudy/api/repository"
)

type Applications struct {
  repository repository.Repositorier
}

func NewApplications(session repository.Repositorier) *Applications {
  return &Applications{repository: session}
}

func (a *Applications) Post(values interface{}) error {
  return nil
}

func (a *Applications) Get(filters interface{}) (interface{}, error) {
  var (
    entity  = []datamodels.ApplicationsFields{}
    err	    error
  )

  if _, err = a.repository.Read(filters, &entity); err != nil {
    return entity, err
  }

  return entity, err
}

func (a *Applications) Delete(conditions interface{}) error {
  return nil
}

func (a *Applications) Put(fields, data interface{}) error {
  return nil
}

/*func (a *Applications) Get(filters interface{}) (i interface{}, err error) {
  return i, err
}

func (a *Applications) Post(values interface{}) error {
  var (
    applications = values.(*datamodels.ApplicationsFields)
  )

  fmt.Println(applications)
  return nil
}*/
