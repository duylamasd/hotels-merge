package integration_test

import (
	"testing"

	"github.com/duylamasd/hotels-merge/bootstrap"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx/fxtest"
)

func TestApp(t *testing.T) {
	app := fxtest.New(t, bootstrap.Modules)

	app.RequireStart()
	app.RequireStop()
}
