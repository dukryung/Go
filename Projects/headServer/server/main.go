package main

import (
	"os"
	"os/signal"
	"personal/Go/Projects/headServer/server/app"
	"personal/Go/Projects/headServer/types/config"
	"syscall"
)

func main() {

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	app := app.NewApp(config.DefaultConfig)

	err := app.RunServers()
	if err != nil {
		panic(err)
	}

	<- quit
}
