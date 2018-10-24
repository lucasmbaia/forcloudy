package services

import (
  "github.com/lucasmbaia/forcloudy/api/config"
  "github.com/lucasmbaia/forcloudy/api/datamodels"
  "github.com/lucasmbaia/forcloudy/api/models"
  "github.com/lucasmbaia/forcloudy/api/repository"
)

type CustomersService interface {
  ResourceService
}

func NewCustomersService() CustomersService {
  return &resourceService{
    fields:	func() interface{} {
      return &datamodels.CustomersFields{}
    },
    model:	func(r repository.Repositorier) models.Models {
      return models.NewCustomers(r)
    },
    repository: config.EnvSingleton.DBConnection,
  }
}
