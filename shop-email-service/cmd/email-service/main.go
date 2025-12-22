package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hawful70/shop-email-service/internal/config"
	"github.com/hawful70/shop-email-service/internal/email"
	kafkamq "github.com/hawful70/shop-email-service/internal/messaging/kafka"
)

func main() {
	cfg := config.MustLoad()
	logger := log.New(os.Stdout, "email-service ", log.LstdFlags)

	if cfg.SMTPHost == "" {
		logger.Fatal("SMTP_HOST is required")
	}
	logger.Printf("using SMTP mailer via %s:%d (TLS=%t)", cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUseTLS)
	mailer := email.NewSMTPMailer(logger, email.SMTPConfig{
		Host:     cfg.SMTPHost,
		Port:     cfg.SMTPPort,
		Username: cfg.SMTPUsername,
		Password: cfg.SMTPPassword,
		From:     cfg.MailFrom,
		FromName: cfg.MailFromName,
		UseTLS:   cfg.SMTPUseTLS,
	})
	handler := email.NewUserCreatedHandler(mailer)
	if len(cfg.KafkaBrokers) == 0 {
		logger.Fatal("no KAFKA_BROKERS configured")
	}
	consumer := kafkamq.NewConsumer(cfg.KafkaBrokers, cfg.KafkaUserCreatedTopic, cfg.KafkaGroupID)
	defer consumer.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Printf("starting email consumer with %d workers\n", cfg.WorkerCount)
	msgCh, errCh := consumer.Stream(ctx)

	var wg sync.WaitGroup
	for i := 0; i < cfg.WorkerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for msg := range msgCh {
				if err := handler.Handle(ctx, msg.Value); err != nil {
					logger.Printf("worker %d failed to handle message: %v", id, err)
				}
				if err := consumer.Commit(ctx, msg); err != nil {
					logger.Printf("worker %d failed to commit: %v", id, err)
				}
			}
		}(i + 1)
	}

	var streamErr error
	select {
	case <-ctx.Done():
	case streamErr = <-errCh:
	}

	wg.Wait()

	if streamErr != nil {
		logger.Fatalf("consumer error: %v", streamErr)
	}

	logger.Println("email service stopped gracefully")
}
