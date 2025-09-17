package e2e_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/duylamasd/hotels-merge/api/domains"
	"github.com/duylamasd/hotels-merge/bootstrap"
	"github.com/duylamasd/hotels-merge/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type TestApp struct {
	Router *gin.Engine
	Server *httptest.Server
}

func NewTestApp(router *gin.Engine) *TestApp {
	return &TestApp{
		Router: router,
		Server: httptest.NewServer(router),
	}
}

func setupTestApp(t *testing.T) (*TestApp, func()) {
	var testApp *TestApp

	app := fxtest.New(t, bootstrap.Modules, fx.Provide(NewTestApp), fx.Populate(&testApp))

	startCtx, cancel := context.WithTimeout(context.Background(), app.StartTimeout())
	defer cancel()

	err := app.Start(startCtx)
	require.NoError(t, err)

	cleanup := func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), app.StopTimeout())
		defer cancel()
		app.Stop(stopCtx)

		if testApp.Server != nil {
			testApp.Server.Close()
		}
	}

	return testApp, cleanup
}

func TestGetV1Hotels(t *testing.T) {
	testApp, cleanup := setupTestApp(t)
	defer cleanup()

	t.Run("GET /api/v1/hotels returns 400 with no query params", func(t *testing.T) {
		resp, err := http.Get(testApp.Server.URL + "/api/v1/hotels")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body domains.HttpError
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, body.Code)
		assert.Equal(t, body.Message, "Either destination or list of hotel ids need to be provided")
	})

	t.Run("GET /api/v1/hotels returns 200 with destination_id", func(t *testing.T) {
		resp, err := http.Get(testApp.Server.URL + "/api/v1/hotels?destination_id=dest1")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var hotels []*sqlc.Hotel
		err = json.NewDecoder(resp.Body).Decode(&hotels)
		assert.NoError(t, err)
	})

	t.Run("GET /api/v1/hotels returns 200 with hotel_ids", func(t *testing.T) {
		resp, err := http.Get(testApp.Server.URL + "/api/v1/hotels?hotel_ids=hotel1&hotel_ids=hotel2")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var hotels []*sqlc.Hotel
		err = json.NewDecoder(resp.Body).Decode(&hotels)
		assert.NoError(t, err)
	})

	t.Run("GET /api/v1/hotels returns 400 with invalid hotel_ids", func(t *testing.T) {
		resp, err := http.Get(testApp.Server.URL + "/api/v1/hotels?hotel_ids=")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body domains.HttpError
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, body.Code)
	})

	t.Run("GET /api/v1/hotels returns 400 with invalid destination_id", func(t *testing.T) {
		resp, err := http.Get(testApp.Server.URL + "/api/v1/hotels?destination_id=")
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body domains.HttpError
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, body.Code)
	})
}
