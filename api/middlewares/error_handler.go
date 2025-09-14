package middlewares

import (
	"net/http"

	"github.com/duylamasd/hotels-merge/api/domains"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ErrorHandler struct {
	logger *zap.Logger
}

func (h *ErrorHandler) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case domains.HttpError:
				c.AbortWithStatusJSON(e.Code, e)
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, domains.NewHttpError(
					http.StatusInternalServerError,
					"Unexpected error occurred",
				))
			}
		}
	}
}

func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{logger: logger}
}
