package media

import (
	"google.golang.org/grpc"
	"net"
	"personal/Go/Projects/video_backend/server/types"
)

type Server struct {
	GrpcServer *grpc.Server
	config types.AppConfig

	close chan bool
}

func NewServer(config types.AppConfig) *Server {
	return &Server{
		config: config,
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