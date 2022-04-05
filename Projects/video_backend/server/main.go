package main

import (
	"os"
	"os/signal"
	"personal/Go/Projects/video_backend/server/app"
	"syscall"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	app := app.NewApp()

	err := app.RunServers()
	if err != nil {
		panic(err)
	}

	<- quit

	app.CloseServers()

}
