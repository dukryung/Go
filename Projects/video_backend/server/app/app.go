package app

import (
	"personal/Go/Projects/video_backend/server/media"
	"personal/Go/Projects/video_backend/server/types"
)

type App struct {
	servers []types.Server
	appConfig types.AppConfig
}

func NewApp(configPath string) *App {
	app := App{}

	app.appConfig = types.AppConfig{}
	err := app.appConfig.LoadAppConfig(configPath)
	if err != nil {
		panic(err)
	}

	app.initServers()

	return &app
}

func (app *App) initServers() {
	mediaServer := media.NewServer(app.appConfig)

	app.servers = append(app.servers,mediaServer)

}

func (app *App) RunServers() error {
	for _, server := range app.servers {
		go server.Run()
	}
	return nil
}

func (app *App) CloseServers() {
	for _, server := range app.servers {
		server.Close()
	}
}
