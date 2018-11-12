package main

import (
	"context"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/lucasmbaia/forcloudy/api/config"
	"github.com/lucasmbaia/forcloudy/api/controllers"
	"github.com/lucasmbaia/forcloudy/api/core-xmpp"
	"github.com/lucasmbaia/forcloudy/api/repository/gorm"
	"github.com/lucasmbaia/forcloudy/api/services"
)

func main() {
	config.EnvDB = config.Database{
		gorm.Config{
			Username:     "forcloudy",
			Password:     "123456",
			Host:         "localhost",
			Port:         "3306",
			DBName:       "forcloudy",
			Timeout:      "10000ms",
			Debug:        true,
			ConnsMaxIdle: 10,
			ConnsMaxOpen: 10,
		},
	}

	config.EnvXmpp = config.Xmpp{
		Host: "192.168.204.131",
		Port: "5222",
		MechanismAuthenticate: "PLAIN",
		User:     "zeus@localhost",
		Password: "totvs@123",
		Room:     "minions@conference.localhost",
	}

	config.LoadConfig()

	go func() {
		if err := core.Run(context.Background()); err != nil {
			panic(err)
		}
	}()

	app := iris.New()

	mvc.Configure(app.Party("/customers"), customers)
	mvc.Configure(app.Party("/customers/{Customer:string}/applications"), applications)

	app.Run(
		// Start the web server at localhost:8080
		iris.Addr("localhost:8080"),
		// disables updates:
		iris.WithoutVersionChecker,
		// skip err server closed when CTRL/CMD+C pressed:
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)
}

func customers(app *mvc.Application) {
	app.Register(services.NewCustomersService())
	app.Handle(new(controllers.CustomerController))
	//iris.New().Party("/{customer:string}/applications").Configure(applications)
}

func applications(app *mvc.Application) {
	app.Register(services.NewApplicationService())
	app.Handle(new(controllers.ApplicationController))
}
