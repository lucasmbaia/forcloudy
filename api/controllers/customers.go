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

func (c *CustomerController) Get() ([]datamodels.CustomersFields, error) {
  var (
    results interface{}
    err	    error
  )

  results, err = c.Services.Get()
  return results.([]datamodels.CustomersFields), err
}

func (c *CustomerController) GetBy(id string) (customer datamodels.CustomersFields, found bool) {
  return datasource.Customers[1], true
}

func (c *CustomerController) Post() error {
  return c.Services.Post()
}
