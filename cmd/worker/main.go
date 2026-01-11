package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"github.com/justshev/go-micro-template/internal/config"
	"github.com/justshev/go-micro-template/internal/logger"
)

func main() {
	cfg, err := config.Load("worker-service")
	if err != nil {
		fmt.Println("config error:", err)
		os.Exit(1)
	}

	log := logger.New(cfg.ServiceName, cfg.LogLevel)
	log.Info("starting worker...")

	amqpURL := cfg.AMQPURL
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		log.Error("failed to connect to rabbitmq: " + err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Error("failed to open channel: " + err.Error())
		os.Exit(1)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"task.created.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("failed to declare queue: " + err.Error())
		os.Exit(1)
	}

	msgs, err := ch.Consume(
		"task.created.queue",
		"",
		false, // auto-ack = false (kita ack manual)
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("failed to consume: " + err.Error())
		os.Exit(1)
	}

	// graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	log.Info("worker is consuming task.created.queue")

	for {
		select {
		case <-done:
			log.Info("shutdown signal received, exiting")
			return
		case msg := <-msgs:
			// kalau channel ditutup, msg.Body bisa kosong; handle sederhana:
			if msg.Body == nil {
				continue
			}

			log.Info("received message: " + string(msg.Body))

			// simulate sending email / notif
			time.Sleep(1 * time.Second)
			log.Info("email mock sent")

			// ACK supaya message dianggap selesai
			_ = msg.Ack(false)
		}
	}
}
