package main

import (
	"Go/Projects/SideProject/common"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("[LOG] starting server...")

	databasename := "sideproject"

	router := common.MakeHandler(databasename)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err == http.ErrServerClosed {
			log.Fatalf("[ERR] Failed to initialize server: %v\n", err)
		}
	}()

	log.Printf("[LOG]Listening on port %v\n", server.Addr)

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERR] server forced to shutdown : %v\n", err)
	}

}
