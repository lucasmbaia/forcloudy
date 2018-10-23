package gorm

import (
  "fmt"
  _ "github.com/go-sql-driver/mysql"
  _gorm "github.com/jinzhu/gorm"
  "time"
  "reflect"
  "errors"
)

const (
  DBTemplate = "%s:%s@tcp(%s:%s)/%s?parseTime=true&timeout=%s"
)

type Client struct {
  DB *_gorm.DB
}

type Config struct {
  Username         string `json:",omitempty"`
  Password         string `json:",omitempty"`
  Host             string `json:",omitempty"`
  Port             string `json:",omitempty"`
  DBName           string `json:",omitempty"`
  Timeout          string `json:",omitempty"`
  Debug            bool   `json:",omitempty"`
  ConnsMaxIdle     int    `json:",omitempty"`
  ConnsMaxOpen     int    `json:",omitempty"`
  ConnsMaxLifetime int    `json:",omitempty"`
}

func NewClient(c Config) (*Client, error) {
  var (
    cli = &Client{}
    err error
  )

  if cli.DB, err = _gorm.Open("mysql", fmt.Sprintf(DBTemplate, c.Username, c.Password, c.Host, c.Port, c.DBName, c.Timeout)); err != nil {
    return cli, err
  }

  cli.DB.LogMode(c.Debug)
  cli.DB.DB().SetMaxIdleConns(c.ConnsMaxIdle)
  cli.DB.DB().SetMaxOpenConns(c.ConnsMaxOpen)
  cli.DB.DB().SetConnMaxLifetime(time.Duration(c.ConnsMaxLifetime))

  return cli, nil
}

func (c *Client) Create(entity interface{}) error {
  if err := c.DB.Create(entity).Error; err != nil {
    return err
  }

  return nil
}

func (c *Client) Delete(condition interface{}) error {
  return nil
}

func (c *Client) Read(condition, entity interface{}) (bool, error) {
  if reflect.ValueOf(entity).Kind() != reflect.Ptr {
    return false, errors.New("The target struct is required to be a pointer")
  }

  var operation *_gorm.DB = c.DB.Set("gorm:auto_preload", true).Find(entity, condition)

  if operation.RecordNotFound() {
    return false, nil
  }

  if operation.Error != nil {
    return true, operation.Error
  }

  return true, nil
}

func (c *Client) Update(condition, entity interface{}) error {
  return nil
}
