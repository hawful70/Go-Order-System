package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/hawful70/shop-identity-service/internal/config"
	"github.com/hawful70/shop-identity-service/internal/httpserver"
	"github.com/hawful70/shop-identity-service/internal/identity"
	"github.com/hawful70/shop-identity-service/internal/identity/domain"
	"github.com/hawful70/shop-identity-service/internal/identity/repository"
	identityhttp "github.com/hawful70/shop-identity-service/internal/identity/transport/http"
)

func main() {
	cfg := config.MustLoad()

	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&domain.UserModel{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	jwtManager := identity.NewJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTExpiresIn)
	repo := repository.NewPostgresRepository(db)
	svc := identity.NewService(repo, jwtManager)
	h := identityhttp.NewHandler(svc)

	r := chi.NewRouter()

	// Basic healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth routes
	r.Route("/api/v1", func(r chi.Router) {
		h.RegisterRoutes(r)
	})

	srv := httpserver.New(":"+cfg.HTTPPort, r)

	go func() {
		if err := srv.Start(); err != nil && err.Error() != "http: Server closed" {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}
	log.Println("identity service stopped gracefully")
}
