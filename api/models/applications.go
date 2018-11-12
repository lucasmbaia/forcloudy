package models

import (
	"errors"
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

func (a *Applications) Post(values interface{}) error {
	var (
		application  = values.(*datamodels.ApplicationsFields)
		applications interface{}
		err          error
	)

	if applications, err = a.Get(datamodels.ApplicationsFields{Name: application.Name, Customer: application.Customer}); err != nil {
		return err
	}

	if len(applications.([]datamodels.ApplicationsFields)) > 0 {
		return errors.New(fmt.Sprintf("Name of application %s exists in database", application.Name))
	}

	if err = a.repository.Create(application); err != nil {
		return err
	}

	return nil
}

func (a *Applications) Get(filters interface{}) (interface{}, error) {
	var (
		entity = []datamodels.ApplicationsFields{}
		err    error
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
