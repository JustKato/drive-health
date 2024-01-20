package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tea.chunkbyte.com/kato/drive-health/lib/config"
	"tea.chunkbyte.com/kato/drive-health/lib/svc"
	"tea.chunkbyte.com/kato/drive-health/lib/web"
)

func main() {
	// Init the database
	svc.InitDB()

	router := web.SetupRouter()

	srv := &http.Server{
		Addr:    config.GetConfiguration().Listen,
		Handler: router,
	}

	// Run the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Run the hardware service
	svc.RunLoggerService()

	// Setting up signal capturing
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
