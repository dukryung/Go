package app

import (
	"fmt"
	"personal/Go/Projects/headServer/server/test"
	"personal/Go/Projects/headServer/types"
	"personal/Go/Projects/headServer/types/config"
)

type App struct {
	config config.AppConfig
	servers []types.Server
}



func NewApp(configPath string) *App {
	app := App{}
	app.config = config.AppConfig{}

	err := app.config.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	app.initServers()


	return &app
}


func (app *App) initServers() {
	server, err :=  test.NewServer(app.config.Test)
	if err != nil {
		panic(err)
	}

	app.servers = append(app.servers, server)

}

func (app *App) RunServers() error {
	if len(app.servers) == 0 {
		return fmt.Errorf("no server enabled")
	}

	for _, server := range app.servers {
		go server.Run()
	}

	return nil
}

