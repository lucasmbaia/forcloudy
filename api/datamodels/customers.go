package datamodels

import (
  "time"
)

type CustomersFields struct {
  ID           string               `json:"id,omitempty"`
  Name         string               `json:"name,omitempty"`
  Applications []ApplicationsFields `json:"applications,omitempty" gorm:"association_foreignkey:ID;foreignkey:Customer"`
  CreatedAt    time.Time            `json:"created_at,omitempty"`
}

func(CustomersFields) TableName() string {
  return "customers"
}
