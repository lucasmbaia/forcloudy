package controllers

import (
	"github.com/kataras/iris"
	"github.com/lucasmbaia/forcloudy/api/datamodels"
	"github.com/lucasmbaia/forcloudy/api/services"
)

type ApplicationController struct {
	ResourceController

	Service services.ResourceService
	ctx     iris.Context
}

/*func (a *ApplicationController) BeginRequest(ctx iris.Context) {
	  a.Service.Set("Customer", ctx.Params().Get("Customer"))
  }

  func (a *ApplicationController) EndRequest(ctx iris.Context) {
  }*/

func (a *ApplicationController) Get() (results []datamodels.ApplicationsFields) {
	return results
	//return a.Service.Get()
}

func (a *ApplicationController) GetByCustomer(customer string) (results []datamodels.ApplicationsFields) {
	return results
}
