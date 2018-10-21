package models

import (
	"fmt"
	//"github.com/lucasmbaia/forcloudy/api/config"
	"github.com/lucasmbaia/forcloudy/api/datamodels"
	"github.com/lucasmbaia/forcloudy/api/repository"
)

type Customers struct {
	repository repository.Repositorier
}

func NewCustomers(session interface{}) *Customers {
	return &Customers{}
}

func (c *Customers) Get() (i interface{}, err error) {
	return i, err
}

func (c *Customers) Post(values interface{}) {
	var (
		customers = values.(*datamodels.CustomersFields)
	)

	fmt.Println(customers)
}
