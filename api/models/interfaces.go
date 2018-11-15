package models

type Models interface {
	Post(interface{}) error
	Get(interface{}) (interface{}, error)
	Delete(interface{}) error
	Put(interface{}, interface{}) error
	Patch(interface{}, interface{}) error
}
