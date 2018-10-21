package models

import (
	"fmt"
	"github.com/lucasmbaia/forcloudy/api/datamodels"
)

type Applications struct {
}

func NewApplications() *Applications {
	return &Applications{}
}

func (a *Applications) Post(values interface{}) {
	var (
		applications = values.(*datamodels.ApplicationsFields)
	)

	fmt.Println(applications)
}
