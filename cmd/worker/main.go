package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justshev/go-micro-template/internal/config"
	"github.com/justshev/go-micro-template/internal/logger"
)

func main() {
	cfg, err := config.Load("worker-service")
	if err != nil {
		fmt.Println("config error:", err)
		os.Exit(1)
	}

	log := logger.New(cfg.ServiceName)
	log.Info("starting...")

	// Simulate worker loop (nanti diganti consume RabbitMQ di Minggu 4)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-done:
			log.Info("shutdown signal received")
			log.Info("stopped cleanly")
			return
		case <-ticker.C:
			log.Info("worker heartbeat")
		}
	}
}
