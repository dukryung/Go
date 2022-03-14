package test

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"personal/Go/Projects/headServer/types/config"
)

type Server struct {
	router *mux.Router
	srv    *http.Server

	close chan bool
}

func NewServer(config config.ServerConfig) (*Server, error){
	s := Server{
		router:mux.NewRouter(),
		close: make(chan bool),
	}
	s.srv = &http.Server{
		Handler: mux.NewRouter(),
		Addr: fmt.Sprintf(":%v", config.Port),
	}

	return &s, nil
}


func (s *Server) Run() {

	s.run()
	<- s.close

}

func (s *Server) run () {
	go func(){
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