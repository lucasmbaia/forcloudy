package config

import (
  "github.com/lucasmbaia/forcloudy/api/repository/gorm"
)

func loadDB() {
  var err error

  if EnvSingleton.DBConnection, err = gorm.NewClient(gorm.Config{
    Username:	      EnvDB.Username,
    Password:	      EnvDB.Password,
    Host:	      EnvDB.Host,
    Port:	      EnvDB.Port,
    DBName:	      EnvDB.DBName,
    Timeout:	      EnvDB.Timeout,
    Debug:	      EnvDB.Debug,
    ConnsMaxIdle:     EnvDB.ConnsMaxIdle,
    ConnsMaxOpen:     EnvDB.ConnsMaxOpen,
    ConnsMaxLifetime: EnvDB.ConnsMaxLifetime,
  }); err != nil {
    panic(err)
  }
}
