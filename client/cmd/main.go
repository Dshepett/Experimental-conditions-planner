package main

import (
	"client/internal/app"
	"client/internal/config"
	"fmt"
)

func main() {
	config := config.New(".env")
	fmt.Print(config)
	app := app.NewApp(config)
	app.Run()
}
