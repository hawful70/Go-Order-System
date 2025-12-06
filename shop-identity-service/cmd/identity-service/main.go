package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/hawful70/shop-identity-service/internal/config"
	"github.com/hawful70/shop-identity-service/internal/httpserver"
	"github.com/hawful70/shop-identity-service/internal/identity"
	identityhttp "github.com/hawful70/shop-identity-service/internal/identity/transport/http"
)

// simple in-memory repository for now
type inMemoryRepo struct {
	mu    sync.RWMutex
	users map[string]identity.User // key: email
}

func newInMemoryRepo() identity.Repository {
	return &inMemoryRepo{
		users: make(map[string]identity.User),
	}
}

func (r *inMemoryRepo) CreateUser(ctx context.Context, u identity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.Email] = u
	return nil
}

func (r *inMemoryRepo) GetUserByEmail(ctx context.Context, email string) (identity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[email]
	if !ok {
		return identity.User{}, identity.ErrInvalidLogin // reuse login-style error
	}
	return u, nil
}

func main() {
	cfg := config.MustLoad()

	jwtManager := identity.NewJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTExpiresIn)
	repo := newInMemoryRepo()
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
