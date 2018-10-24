package datamodels

import (
  "time"
)

type ApplicationsFields struct {
  ID		  string    `json:"id,omitempty"`
  Customer	  string    `json:"customer,omitempty"`
  Name		  string    `json:"name,omitempty"`
  Description	  string    `json:"description,omitempty"`
  Cpus		  string    `json:"cpus,omitempty"`
  Memory	  int	    `json:"memory,omitempty"`
  TotalContainers int	    `json:"totalContainers,omitempty"`
  Ports		  []Ports   `json:"ports,omitempty" gorm:"association_foreignkey:ID;foreignkey:Application"`
  CreatedAt	  time.Time `json:"created_at,omitempty"`
}

type Ports struct {
  Application string  `json:"-"`
  Port	      int     `json:"port,omitempty"`
  Protocol    string  `json:"protocol,omitempty"`
}

func(ApplicationsFields) TableName() string {
  return "applications"
}

func(Ports) TableName() string {
  return "ports"
}
