package models

type Models interface {
  Get(interface{}) (interface{}, error)
  Post(interface{}) error
}
