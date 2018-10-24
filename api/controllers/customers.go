package controllers

import (
  "github.com/lucasmbaia/forcloudy/api/datamodels"
  "github.com/lucasmbaia/forcloudy/api/datasource"
)

type CustomerController struct {
  ResourceController
}

/*func (c *CustomerController) Get() (results []datamodels.CustomersFields) {
  var data []datamodels.CustomersFields

  for _, v := range datasource.Customers {
    data = append(data, v)
  }

  //c.Ctx.ContentType("application/json")
  return data
}*/

func (c *CustomerController) Get() (data []datamodels.CustomersFields, err error) {
  /*var (
    results interface{}
    err	    error
  )

  //results, err = c.Services.Get()
  return results.([]datamodels.CustomersFields), err*/
  c.Services.Get(c.Ctx)
  return data, err
}

func (c *CustomerController) GetBy(id string) (customer datamodels.CustomersFields, found bool) {
  return datasource.Customers[1], true
}

func (c *CustomerController) Post() (response datamodels.Response, error) {
  return c.Services.Post(c.Ctx)
}
