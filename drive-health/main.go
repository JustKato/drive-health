package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JustKato/drive-health/lib/config"
	"github.com/JustKato/drive-health/lib/svc"
	"github.com/JustKato/drive-health/lib/web"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("[🟨] No .env file found")
	}

	// Init the database
	svc.InitDB()
	cfg := config.GetConfiguration()

	router := web.SetupRouter()

	srv := &http.Server{
		Addr:    cfg.Listen,
		Handler: router,
	}

	// Run the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[🛑] listening failed: %s\n", err)
		}
	}()

	// Run the hardware service
	svc.RunLoggerService()
	// Run the cleanup service
	svc.RunCleanupService()

	// Setting up signal capturing
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
	<-quit
	log.Println("[🦝] Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("[🛑] Server forced to shutdown:", err)
	}

	log.Println("[🦝] Server exiting")
}
