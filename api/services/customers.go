package services

import (
  "github.com/lucasmbaia/forcloudy/api/config"
  "github.com/lucasmbaia/forcloudy/api/datamodels"
  "github.com/lucasmbaia/forcloudy/api/models"
)

type CustomersService interface {
  ResourceService
}

func NewCustomersService() CustomersService {
  return &resourceService{
    fields:     &datamodels.CustomersFields{},
    model:      models.NewCustomers(config.EnvSingleton.DBConnection),
    repository: config.EnvSingleton.DBConnection,
  }
}
