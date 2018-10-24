package services

import (
  "github.com/lucasmbaia/forcloudy/api/config"
  "github.com/lucasmbaia/forcloudy/api/datamodels"
  "github.com/lucasmbaia/forcloudy/api/models"
  "github.com/lucasmbaia/forcloudy/api/repository"
)

type ApplicationsService interface {
  ResourceService
}

func NewApplicationService() CustomersService {
  return &resourceService{
    fields:	func() interface{} {
      return &datamodels.ApplicationsFields{}
    },
    model:	func(r repository.Repositorier) models.Models {
      return models.NewApplications(r)
    },
    repository: config.EnvSingleton.DBConnection,
  }
}
