package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/hawful70/shop-identity-service/internal/identity"
)

type Handler struct {
	svc        identity.Service
	jwtManager *identity.JWTManager
}

func NewHandler(svc identity.Service, jwtManager *identity.JWTManager) *Handler {
	return &Handler{svc: svc, jwtManager: jwtManager}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	// Route registration logic goes here
	r.Post("/auth/register", h.handleRegister)
	r.Post("/auth/login", h.handleLogin)

	r.Group(func(protected chi.Router) {
		protected.Use(h.jwtAuthMiddleware)
		protected.Get("/auth/me", h.handleMe)
	})
}

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Register(r.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		switch err {
		case identity.ErrEmailTaken:
			http.Error(w, err.Error(), http.StatusConflict)
		case identity.ErrPasswordTooWeak, identity.ErrEmailRequired:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	res := registerResponse{
		ID:       string(user.ID),
		Email:    user.Email,
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(res)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	_, token, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if err == identity.ErrInvalidLogin {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	res := loginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

type meResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (h *Handler) handleMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := identity.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	res := meResponse{
		ID:       claims.UserID,
		Email:    claims.Email,
		Username: claims.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (h *Handler) jwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
			http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := h.jwtManager.VerifyToken(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := identity.ContextWithClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
