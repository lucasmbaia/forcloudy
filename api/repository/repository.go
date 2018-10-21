package repository

import (
	"github.com/lucasmbaia/forcloudy/api/repository/gorm"
)

type Repositorier interface {
	Create(intity interface{}) error
	Delete(condition interface{}) error
	Read(condition, entity interface{}) (bool, error)
	Update(condition, entity interface{}) error
}

func NewGormRepository(c gorm.Config) Repositorier {
	var cli *gorm.Client

	cli, _ = gorm.NewClient(c)

	return cli
}
