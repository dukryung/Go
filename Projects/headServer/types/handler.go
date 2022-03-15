package types

import "github.com/gorilla/mux"

type Handler interface {
	Run()
	RegisterRoute(route *mux.Router) *mux.Router
}


type HandlerManager struct {
	handlers []Handler
}

func NewHandlerManager(handlers ...Handler) *HandlerManager {
	return &HandlerManager{
		handlers: handlers,
	}
}

func (m *HandlerManager) Run() {
	for _, handler := range m.handlers {
		go handler.Run()
	}
}

func (m *HandlerManager) registerRoute(router *mux.Router) *mux.Router {
	for _, handler := range  m.handlers {
		handler.RegisterRoute(router)
	}
	return router
}
