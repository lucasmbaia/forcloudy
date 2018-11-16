package config

import (
	"github.com/lucasmbaia/forcloudy/logging"
)

func LoadLog(level string) {
	EnvSingleton.Log = logging.New(level)
}
