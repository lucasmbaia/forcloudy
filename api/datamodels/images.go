package datamodels

import (
	"time"
)

type ImagesFields struct {
	ID        string    `json:"id,omitempty"`
	Customer  string    `json:"customer,omitempty"`
	Name      string    `json:"name,omitempty"`
	Version   string    `json:"version,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (ImagesFields) TableName() string {
	return "images"
}
