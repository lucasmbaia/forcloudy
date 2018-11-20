package logging

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"reflect"
)

const (
	WARN  = "WARN"
	DEBUG = "DEBUG"
	INFO  = "INFO"
)

type Logger struct {
	cli *logrus.Logger
}

func New(level string) *Logger {
	var log = logrus.New()
	log.Level = logLevel(level)
	log.Formatter = new(logrus.TextFormatter)

	return &Logger{cli: log}
}

func logLevel(level string) logrus.Level {
	switch level {
	case WARN:
		return logrus.WarnLevel
	case DEBUG:
		return logrus.DebugLevel
	default:
		return logrus.InfoLevel
	}
}

func (l *Logger) Infofc(template string, args ...interface{}) {
	l.cli.Infof(template, setArgs(args...)...)
}

func (l *Logger) Debugfc(template string, args ...interface{}) {
	l.cli.Debugf(template, setArgs(args...)...)
}

func (l *Logger) Errorfc(template string, args ...interface{}) {
	l.cli.Errorf(template, setArgs(args...)...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.cli.Infof(template, args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.cli.Debugf(template, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.cli.Errorf(template, args...)
}

func setArgs(args ...interface{}) []interface{} {
	var (
		fields  []interface{}
		element reflect.Value
		body    []byte
		err     error
	)

	for _, arg := range args {
		element = reflect.Indirect(reflect.ValueOf(arg))

		if element.Kind() == reflect.Ptr {
			element = element.Elem()
		}

		if element.Kind() == reflect.Struct || element.Kind() == reflect.Slice {
			if body, err = json.Marshal(arg); err != nil {
				fields = append(fields, arg)
			} else {
				fields = append(fields, string(body))
			}

			continue
		}

		fields = append(fields, arg)
	}

	return fields
}
