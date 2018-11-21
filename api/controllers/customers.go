package controllers

import (
	"github.com/lucasmbaia/forcloudy/api/datamodels"
)

type CustomerController struct {
	ResourceController
}

func (c *CustomerController) Get() ([]datamodels.CustomersFields, error) {
	var (
		results interface{}
		err     error
	)

	results, err = c.Services.Get(c.Ctx)
	return results.([]datamodels.CustomersFields), err
}

func (c *CustomerController) GetBy(id string) ([]datamodels.CustomersFields, error) {
	var (
		results interface{}
		err     error
	)

	results, err = c.Services.GetById(c.Ctx, id)
	return results.([]datamodels.CustomersFields), err
}

func (c *CustomerController) Post() (datamodels.Response, error) {
	return c.Services.Post(c.Ctx)
}
