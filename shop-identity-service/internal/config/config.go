package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort              string
	GRPCPort              string
	JWTSecret             string
	JWTIssuer             string
	JWTExpiresIn          time.Duration
	DBDSN                 string
	KafkaBrokers          []string
	KafkaUserCreatedTopic string
}

func Load() Config {
	_ = godotenv.Load()

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8081" // identity service default
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "9091"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: JWT_SECRET is not set, using insecure default (do not use in production)")
		secret = "dev-insecure-secret"
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "shop-identity-service"
	}

	expStr := os.Getenv("JWT_EXPIRES_IN") // e.g. "15m" or "1h"
	if expStr == "" {
		expStr = "15m"
	}
	exp, err := time.ParseDuration(expStr)
	if err != nil {
		log.Printf("invalid JWT_EXPIRES_IN=%s, fallback to 15m\n", expStr)
		exp = 15 * time.Minute
	}

	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		dbDSN = "postgres://postgres:postgres@localhost:5432/identity?sslmode=disable"
	}

	kafkaBrokersEnv := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokersEnv == "" {
		kafkaBrokersEnv = "localhost:19092"
	}
	var kafkaBrokers []string
	for _, broker := range strings.Split(kafkaBrokersEnv, ",") {
		broker = strings.TrimSpace(broker)
		if broker != "" {
			kafkaBrokers = append(kafkaBrokers, broker)
		}
	}

	kafkaUserCreatedTopic := os.Getenv("KAFKA_TOPIC_USER_CREATED")
	if kafkaUserCreatedTopic == "" {
		kafkaUserCreatedTopic = "user_created"
	}

	return Config{
		HTTPPort:              httpPort,
		GRPCPort:              grpcPort,
		JWTSecret:             secret,
		JWTIssuer:             issuer,
		JWTExpiresIn:          exp,
		DBDSN:                 dbDSN,
		KafkaBrokers:          kafkaBrokers,
		KafkaUserCreatedTopic: kafkaUserCreatedTopic,
	}
}

func MustLoad() Config {
	return Load()
}
