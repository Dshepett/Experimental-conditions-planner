package main

import (
	"client/internal/app"
	"client/internal/config"
)

func main() {
	config := config.New(".env")
	app := app.NewApp(config)
	app.Run()
}
