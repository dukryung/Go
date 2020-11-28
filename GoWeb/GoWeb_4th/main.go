package main

import (
	"Go/GoWeb/GoWeb_4th/decoHandler"
	"Go/GoWeb/GoWeb_4th/myapp"
	"log"
	"net/http"
	"time"
)

func logger(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Println("[LOGGER1] Started")
	h.ServeHTTP(w, r)
	log.Println("[LOGGER1] Completed time: ", time.Since(start).Milliseconds())
}

func NewHandler() http.Handler {
	mux := myapp.NewHandler()
	h := decoHandler.NewDecoHandler(mux, logger)
	return h

}

func main() {

	mux := NewHandler()
	http.ListenAndServe(":8080", mux)
}
