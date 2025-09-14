package main

import (
	"github.com/duylamasd/hotels-merge/bootstrap"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		bootstrap.Modules,
	).Run()
}
