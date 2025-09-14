package config

import (
	"os"

	"go.uber.org/fx"
)

type Config struct {
	DBUri    string
	Port     string
	Env      string
	LogLevel string
}

func NewConfig() *Config {
	return &Config{
		DBUri:    os.Getenv("DB_URI"),
		Port:     os.Getenv("PORT"),
		Env:      os.Getenv("ENV"),
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}

var Module = fx.Options(
	fx.Provide(NewConfig),
	fx.Provide(NewDBConn),
	fx.Provide(NewDBStore),
)
