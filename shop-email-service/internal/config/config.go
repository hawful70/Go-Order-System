package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBrokers          []string
	KafkaGroupID          string
	KafkaUserCreatedTopic string
	MailFrom              string
	MailFromName          string
	WorkerCount           int
	SMTPHost              string
	SMTPPort              int
	SMTPUsername          string
	SMTPPassword          string
	SMTPUseTLS            bool
}

func Load() Config {
	_ = godotenv.Load()

	brokersEnv := env("KAFKA_BROKERS", "localhost:19092")
	var brokers []string
	for _, b := range strings.Split(brokersEnv, ",") {
		b = strings.TrimSpace(b)
		if b != "" {
			brokers = append(brokers, b)
		}
	}

	groupID := env("KAFKA_GROUP_ID", "email-service")
	topic := env("KAFKA_TOPIC_USER_CREATED", "user_created")
	mailFrom := env("MAIL_FROM", "welcome@example.com")
	mailFromName := env("MAIL_FROM_NAME", "Shop Team")
	workers := envInt("EMAIL_WORKERS", 4)
	smtpHost := env("SMTP_HOST", "")
	smtpPort := envInt("SMTP_PORT", 587)
	smtpUser := env("SMTP_USERNAME", "")
	smtpPass := env("SMTP_PASSWORD", "")
	smtpTLS := envBool("SMTP_USE_TLS", true)

	return Config{
		KafkaBrokers:          brokers,
		KafkaGroupID:          groupID,
		KafkaUserCreatedTopic: topic,
		MailFrom:              mailFrom,
		MailFromName:          mailFromName,
		WorkerCount:           workers,
		SMTPHost:              smtpHost,
		SMTPPort:              smtpPort,
		SMTPUsername:          smtpUser,
		SMTPPassword:          smtpPass,
		SMTPUseTLS:            smtpTLS,
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			return i
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		switch strings.ToLower(v) {
		case "true", "1", "yes", "y":
			return true
		case "false", "0", "no", "n":
			return false
		}
	}
	return fallback
}

func MustLoad() Config {
	return Load()
}
