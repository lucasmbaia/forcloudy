package services

import (
  "github.com/lucasmbaia/forcloudy/api/config"
  "github.com/lucasmbaia/forcloudy/api/datamodels"
  "github.com/lucasmbaia/forcloudy/api/models"
)

type ApplicationsService interface {
  ResourceService
}

func NewApplicationService() CustomersService {
  return &resourceService{
    fields:	&datamodels.ApplicationsFields{},
    model:	models.NewApplications(config.EnvSingleton.DBConnection),
    repository: config.EnvSingleton.DBConnection,
  }
}
