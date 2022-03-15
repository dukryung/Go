package event

import (
	"fmt"
	"github.com/gorilla/mux"
	"personal/Go/Projects/headServer/types"
	"personal/Go/Projects/headServer/types/config"
	"reflect"
)

type Handler struct {
	config config.ServerConfig

}


func NewHandler(ctx types.Context, config config.ServerConfig) *Handler {
	handler := Handler{}
	handler.config = config

	handler.init()

	return &handler
}

func (s *Handler) init() {
	element := reflect.ValueOf(*s)

	for i:= 0; i< element.NumField(); i++ {
		field := element.Field(i)
		fieldType := field.Type()
		fmt.Println("fieldType : ", fieldType)
	}
}

func (s *Handler) Run() {
	fmt.Println("run")
}

func (s *Handler) Close() {
	fmt.Println("close")
}

func (s *Handler) RegisterRoute(router *mux.Router) *mux.Router {
	return router
}