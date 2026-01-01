package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justshev/go-micro-template/internal/config"
	"github.com/justshev/go-micro-template/internal/httpserver"
	"github.com/justshev/go-micro-template/internal/logger"
)

func main() {
	cfg, err := config.Load("api-service")
	if err != nil {
		fmt.Println("config error:", err)
		os.Exit(1)
	}

	log := logger.New(cfg.ServiceName)
	log.Info("starting...")

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	srv := httpserver.New(httpserver.Addr(cfg.HTTPPort), mux)

	// graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("listening on " + httpserver.Addr(cfg.HTTPPort))
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("server error: " + err.Error())
			os.Exit(1)
		}
	}()

	<-done
	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown error: " + err.Error())
		os.Exit(1)
	}
	log.Info("stopped cleanly")
}
