package logger

import (
	"os"
)

type Logger interface {
	Info(msg string, kv ...any)
	Warn(msg string, kv ...any)
	Error(msg string, kv ...any)
	Debug(msg string, kv ...any)
	Infof(template string, args ...any)
	Warnf(template string, args ...any)
	Errorf(template string, args ...any)
	Debugf(template string, args ...any)
}

func NewLogger() Logger {
	if os.Getenv("LOGGER_MODE") == "cute" {
		return NewDevLogger()
	}
	if os.Getenv("APP_ENV") == "production" {
		return NewProdLogger()
	}
	return NewDevLogger()
}
