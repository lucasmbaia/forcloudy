package models

import (
	"github.com/lucasmbaia/forcloudy/api/datamodels"
	"github.com/lucasmbaia/forcloudy/api/repository"
)

type Containers struct {
	repository repository.Repositorier
}

func NewContainers(session repository.Repositorier) *Containers {
	return &Containers{repository: session}
}

func (c *Containers) Post(values interface{}) error {
	var (
		container = values.(*datamodels.ContainersFields)
		err       error
	)

	if err = c.repository.Create(container); err != nil {
		return err
	}

	return nil
}

func (c *Containers) Get(filters interface{}) (interface{}, error) {
	var (
		entity = []datamodels.ContainersFields{}
		err    error
	)

	if _, err = c.repository.Read(filters, &entity); err != nil {
		return entity, err
	}

	return entity, err
}

func (c *Containers) Delete(conditions interface{}) error {
	return nil
}

func (c *Containers) Put(fields, data interface{}) error {
	return nil
}

func (c *Containers) Patch(fields, data interface{}) error {
	var (
		conditions = fields.(*datamodels.ContainersFields)
		entity     = data.(*datamodels.ContainersFields)
	)

	return c.repository.Update(conditions, entity)
}
