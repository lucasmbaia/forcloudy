package gorm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_gorm "github.com/jinzhu/gorm"
	"time"
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
	return nil
}

func (c *Client) Delete(condition interface{}) error {
	return nil
}

func (c *Client) Read(condition, entity interface{}) (bool, error) {
	return false, nil
}

func (c *Client) Update(condition, entity interface{}) error {
	return nil
}
