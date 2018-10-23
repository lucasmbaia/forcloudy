package main

import (
  "github.com/kataras/iris"
  "github.com/kataras/iris/mvc"
  "github.com/lucasmbaia/forcloudy/api/controllers"
  "github.com/lucasmbaia/forcloudy/api/services"
  "github.com/lucasmbaia/forcloudy/api/config"
  "github.com/lucasmbaia/forcloudy/api/repository/gorm"
)

func main() {
  config.EnvDB = config.Database{
    gorm.Config{
      Username:	  "forcloudy",
      Password:	  "123456",
      Host:	  "localhost",
      Port:	  "3306",
      DBName:	  "forcloudy",
      Timeout:	  "10000ms",
      Debug:	  true,
      ConnsMaxIdle: 10,
      ConnsMaxOpen: 10,
    },
  }

  config.LoadConfig()
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
