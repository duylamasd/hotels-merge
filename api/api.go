package api

import (
	v1Controllers "github.com/duylamasd/hotels-merge/api/controllers/v1"
	"github.com/duylamasd/hotels-merge/api/middlewares"
	v1Routes "github.com/duylamasd/hotels-merge/api/routes/v1"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func registerRoutes(
	engine *gin.Engine,
	v1Routes *v1Routes.V1Routes,
	errorHandler *middlewares.ErrorHandler,
) {
	engine.Use(errorHandler.Handler())

	api := engine.Group("/api")
	v1 := api.Group("/v1")
	v1Routes.Register(v1)
}

var Module = fx.Options(
	middlewares.Module,
	v1Controllers.Module,
	v1Routes.Module,
	fx.Invoke(registerRoutes),
)
