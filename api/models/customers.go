package models

import (
	"errors"
	"fmt"
	//"github.com/lucasmbaia/forcloudy/api/config"
	"github.com/lucasmbaia/forcloudy/api/datamodels"
	"github.com/lucasmbaia/forcloudy/api/repository"
)

type Customers struct {
	repository repository.Repositorier
}

func NewCustomers(session repository.Repositorier) *Customers {
	return &Customers{repository: session}
}

func (c *Customers) Post(values interface{}) error {
	var (
		customer  = values.(*datamodels.CustomersFields)
		customers interface{}
		err       error
	)

	if customers, err = c.Get(datamodels.CustomersFields{Name: customer.Name}); err != nil {
		return err
	}

	if len(customers.([]datamodels.CustomersFields)) > 0 {
		return errors.New(fmt.Sprintf("Name of customer %s exists in database", customer.Name))
	}

	if err = c.repository.Create(customer); err != nil {
		return err
	}

	return nil
}

func (c *Customers) Get(filters interface{}) (interface{}, error) {
	var (
		entity = []datamodels.CustomersFields{}
		err    error
	)

	if _, err = c.repository.Read(filters, &entity); err != nil {
		return entity, err
	}

	return entity, err
}

func (c *Customers) Delete(conditions interface{}) error {
	return nil
}

func (c *Customers) Put(fields, data interface{}) error {
	return nil
}

func (c *Customers) Patch(fields, data interface{}) error {
	return nil
}

/*func (c *Customers) Get(filters interface{}) (interface{}, error) {
  var (
    entity  = []datamodels.CustomersFields{}
    err	    error
  )

  if _, err = c.repository.Read(filters, &entity); err != nil {
    return entity, err
  }

  return entity, err
}

func (c *Customers) Post(values interface{}) error {
  var (
    customer  = values.(*datamodels.CustomersFields)
    customers interface{}
    err	      error
  )

  if customers, err = c.Get(datamodels.CustomersFields{Name: customer.Name}); err != nil {
    return nil
  }

  if len(customers.([]datamodels.CustomersFields)) > 0 {
    return errors.New(fmt.Sprintf("Name of customer %s exists in database", customer.Name))
  }

  if err = c.repository.Create(customer); err != nil {
    return err
  }

  return nil
}*/
