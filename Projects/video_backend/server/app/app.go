package app

import "personal/Go/Projects/video_backend/server/types"

type App struct {
	servers []types.Server
}

func NewApp() *App {
	return &App{

	}
}

func (app *App) initServers() {


}


func (app *App) RunServers() {
	for _, server := range app.servers {
		go server.Run()
	}
}

func (app *App) CloseServers() {
	for _, server := range app.servers {
		server.Close()
	}
}
