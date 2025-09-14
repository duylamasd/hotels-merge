package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/duylamasd/hotels-merge/api"
	"github.com/duylamasd/hotels-merge/config"
	"github.com/duylamasd/hotels-merge/lib"
	"github.com/duylamasd/hotels-merge/services"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewGinEngine(logger *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	return r
}

func RegisterHooks(lc fx.Lifecycle, engine *gin.Engine, config *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%s", config.Port)
			go engine.Run(addr)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
}

var Modules = fx.Options(
	config.Module,
	lib.Module,
	fx.Provide(NewGinEngine),
	services.Module,
	api.Module,
	fx.Invoke(RegisterHooks),
)
