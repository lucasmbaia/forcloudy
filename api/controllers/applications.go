package controllers

import (
	"github.com/lucasmbaia/forcloudy/api/datamodels"
)

type ApplicationController struct {
	ResourceController
}

func (a *ApplicationController) Get() ([]datamodels.ApplicationsFields, error) {
	var (
		results interface{}
		err     error
	)

	results, err = a.Services.Get(a.Ctx)
	return results.([]datamodels.ApplicationsFields), err
}

func (a *ApplicationController) GetBy(id string) ([]datamodels.ApplicationsFields, error) {
	var (
		results interface{}
		err     error
	)

	results, err = a.Services.GetById(a.Ctx, id)
	return results.([]datamodels.ApplicationsFields), err
}

func (a *ApplicationController) Post() (datamodels.Response, error) {
	return a.Services.Post(a.Ctx)
}

func (a *ApplicationController) DeleteBy(id string) error {
	return a.Services.DeleteById(a.Ctx, id)
}

/*func (a *ApplicationController) GetByCustomer(customer string) (results []datamodels.ApplicationsFields) {
	return results
}*/
