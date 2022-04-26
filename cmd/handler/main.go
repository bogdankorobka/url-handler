package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bogdankorobka/url-handler/internal/server"
)

func main() {
	log.Println("start application")

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// create server
	s := server.NewServer("localhost:3000")

	// graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, 30*time.Second)
		defer shutdownStopCtx()

		go func() {
			<-shutdownCtx.Done()

			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		s.Stop(shutdownCtx)

		serverStopCtx()

		log.Println("gracefully shutting down")
	}()

	log.Println("start server")
	s.Start()

	// Wait for server context to be stopped
	<-serverCtx.Done()

	log.Println("application shutdowned")
}
