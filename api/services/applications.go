package services

import (
	//"fmt"
	"github.com/lucasmbaia/forcloudy/api/datamodels"
	"github.com/lucasmbaia/forcloudy/api/models"
	//"reflect"
)

type ApplicationsService interface {
	ResourceService
	//Set(field string, value interface{})
	//Get() []datamodels.ApplicationsFields
}

func NewApplicationService() ApplicationsService {
	return &resourceService{
		fields: &datamodels.ApplicationsFields{},
		model:  models.NewApplications(),
	}
	//return &applicationsService{}
}

/*type applicationsService struct {
	model datamodels.ApplicationsFields
}

func (a *applicationsService) Set(field string, value interface{}) {
	reflect.ValueOf(&a.model).Elem().FieldByName(field).Set(reflect.ValueOf(value))
}

func (a *applicationsService) Get() (result []datamodels.ApplicationsFields) {
	fmt.Println(a.model)
	return result
}*/
