package middlewares_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/duylamasd/hotels-merge/api/domains"
	"github.com/duylamasd/hotels-merge/api/middlewares"
	"github.com/duylamasd/hotels-merge/config"
	"github.com/duylamasd/hotels-merge/lib"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestErrorHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger, _ := lib.NewLogger(&config.Config{LogLevel: "info"})
	errorHandler := middlewares.NewErrorHandler(logger)

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.Use(errorHandler.Handler())

	router.GET("/test", func(c *gin.Context) {
		e := domains.NewHttpError(http.StatusBadRequest, "Error")
		_ = c.Error(e)
	})

	t.Run("should return 400 for HttpError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response domains.HttpError
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, response.Code)
		assert.Equal(t, "Error", response.Message)
	})
}
