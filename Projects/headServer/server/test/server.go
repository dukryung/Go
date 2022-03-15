package test

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"personal/Go/Projects/headServer/server/test/event"
	"personal/Go/Projects/headServer/types"
	"personal/Go/Projects/headServer/types/config"
)

var _ types.Server = &Server{}

type Server struct {
	ctx    types.Context
	router *mux.Router
	srv    *http.Server

	eventHandler *event.Handler

	handlerManager *types.HandlerManager

	close chan bool
}

func NewServer(config config.ServerConfig) (*Server, error) {
	s := Server{
		router: mux.NewRouter(),
		close:  make(chan bool),
	}
	s.srv = &http.Server{
		Handler: mux.NewRouter(),
		Addr:    fmt.Sprintf(":%v", config.Port),
	}

	s.eventHandler = event.NewHandler(s.ctx, config)
	s.handlerManager = types.NewHandlerManager(s.eventHandler)

	s.srv = &http.Server{
		Handler: s.router,
		Addr:    fmt.Sprintf("%v", config.Port),
	}

	return &s, nil
}

func (s *Server) Run() {

	s.handlerManager.Run()
	s.run()
	<-s.close

}

func (s *Server) run() {
	go func() {
		err := s.srv.ListenAndServe()
		if err != nil {
			fmt.Println("server listen error")
			panic(err)
		}
	}()
}

func (s *Server) Close() {
	s.close <- true
}
