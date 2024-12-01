package main

import (
	"autoscaler/internal/application"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app, err := application.NewApplication(ctx)
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}

	go func() {
		if err := app.Start(); err != nil {
			log.Printf("Application stopped with error: %v", err)
			cancel()
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
		log.Println("Received shutdown signal")
	case <-ctx.Done():
		log.Println("Context canceled")
	}

	if err = app.Shutdown(ctx); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
}
