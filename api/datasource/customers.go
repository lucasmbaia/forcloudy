package datasource

import (
	"github.com/lucasmbaia/forcloudy/api/datamodels"
	"time"
)

var Customers = map[int]datamodels.CustomersFields{
	1: {
		ID:        "bee2ec29-e6e7-4529-b938-0b90dfc626a7",
		Name:      "lucas",
		CreatedAt: time.Now(),
	},
	2: {
		ID:        "49b65f2d-d1f4-4483-8d04-c16f368d9f5f",
		Name:      "martins",
		CreatedAt: time.Now(),
	},
}
