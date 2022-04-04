package media

import (
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	GrpcServer *grpc.Server

	close chan bool
}

func NewServer() *Server {
	return &Server{

	}
}

func (s *Server) Run(){
	listen, err :=  net.Listen("tcp", ":10333")
	if err != nil {
		panic(err)
	}
	s.GrpcServer.Serve(listen)
}

func (s *Server) Close() {
	s.close <- true
}