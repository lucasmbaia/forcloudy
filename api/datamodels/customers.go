package datamodels

import (
	"time"
)

type CustomersFields struct {
	ID           string               `json:"id,omitemtpy"`
	Name         string               `json:"name,omitemtpy"`
	Applications []ApplicationsFields `json:"applications,omitemtpy" gorm:"association_foreignkey:ID;foreignkey:Customer"`
	CreatedAt    time.Time            `json:"created_at,omitemtpy"`
}
