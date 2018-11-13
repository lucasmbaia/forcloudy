package models

import (
	"errors"
	"fmt"
	"github.com/lucasmbaia/forcloudy/api/datamodels"
	"github.com/lucasmbaia/forcloudy/api/repository"
	"github.com/satori/go.uuid"
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
		customers    interface{}
		image        string
		imageID      uuid.UUID
		containerID  uuid.UUID
		iterator     = 1
	)

	if applications, err = a.Get(datamodels.ApplicationsFields{Name: application.Name, Customer: application.Customer}); err != nil {
		return err
	}

	if len(applications.([]datamodels.ApplicationsFields)) > 0 {
		return errors.New(fmt.Sprintf("Name of application %s exists in database", application.Name))
	}

	if customers, err = NewCustomers(a.repository).Get(
		datamodels.CustomersFields{
			ID: application.Customer,
		},
	); err != nil {
		return err
	} else {
		if len(customers.([]datamodels.CustomersFields)) == 0 {
			return errors.New("Invalid Customer")
		}
	}

	image = fmt.Sprintf("%s_app-%s", customers.([]datamodels.CustomersFields)[0].Name, application.Name)
	if imageID, err = uuid.NewV4(); err != nil {
		return err
	}

	if err = NewImages(a.repository).Post(
		&datamodels.ImagesFields{
			ID:       imageID.String(),
			Customer: application.Customer,
			Name:     image,
			Version:  "v1",
		},
	); err != nil {
		return err
	}

	application.Image = imageID.String()
	application.Status = "IN_PROGRESS"

	if err = a.repository.Create(application); err != nil {
		return err
	}

	for iterator <= application.TotalContainers {
		if containerID, err = uuid.NewV4(); err != nil {
			return err
		}

		if err = NewContainers(a.repository).Post(
			&datamodels.ContainersFields{
				ID:          containerID.String(),
				Customer:    application.Customer,
				Application: application.ID,
				Name:        fmt.Sprintf("%s_app-%s-%d", customers.([]datamodels.CustomersFields)[0].Name, application.Name, iterator),
				Status:      "IN_PROGRESS",
				State:       "CREATING",
			},
		); err != nil {
			return err
		}

		iterator++
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
