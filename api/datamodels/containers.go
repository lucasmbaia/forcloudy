package datamodels

import (
	"time"
)

type ContainersFields struct {
	ID          string    `json:"id,omitempty"`
	Customer    string    `json:"-"`
	Application string    `json:"-"`
	Name        string    `json:"name,omitempty"`
	Status      string    `json:"status,omitempty"`
	State       string    `json:"state,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

func (ContainersFields) TableName() string {
	return "containers"
}
