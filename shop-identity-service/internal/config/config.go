package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort     string
	JWTSecret    string
	JWTIssuer    string
	JWTExpiresIn time.Duration
	DBDSN        string
}

func Load() Config {
	_ = godotenv.Load()

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8081" // identity service default
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

	return Config{
		HTTPPort:     port,
		JWTSecret:    secret,
		JWTIssuer:    issuer,
		JWTExpiresIn: exp,
		DBDSN:        dbDSN,
	}
}

func MustLoad() Config {
	return Load()
}
