package lib

import (
	"github.com/duylamasd/hotels-merge/config"
	"go.uber.org/zap"
)

func NewLogger(config *config.Config) (*zap.Logger, error) {
	var zapConfig zap.Config

	if config.Env == "production" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	level, err := zap.ParseAtomicLevel(config.LogLevel)
	if err != nil {
		return nil, err
	}
	zapConfig.Level = level

	return zapConfig.Build()
}
