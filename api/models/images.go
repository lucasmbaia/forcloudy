package models

import (
	"github.com/lucasmbaia/forcloudy/api/datamodels"
	"github.com/lucasmbaia/forcloudy/api/repository"
)

type Images struct {
	repository repository.Repositorier
}

func NewImages(session repository.Repositorier) *Images {
	return &Images{repository: session}
}

func (i *Images) Post(values interface{}) error {
	var (
		image = values.(*datamodels.ImagesFields)
		err   error
	)

	if err = i.repository.Create(image); err != nil {
		return err
	}

	return nil
}

func (i *Images) Get(filters interface{}) (interface{}, error) {
	var (
		entity = []datamodels.ImagesFields{}
		err    error
	)

	if _, err = i.repository.Read(filters, &entity); err != nil {
		return entity, err
	}

	return entity, err
}

func (i *Images) Delete(conditions interface{}) error {
	return nil
}

func (i *Images) Put(fields, data interface{}) error {
	return nil
}

func (i *Images) Patch(fields, data interface{}) error {
	return nil
}
