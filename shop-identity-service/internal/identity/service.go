package identity

import (
	"context"
	"errors"
	"strings"

	"github.com/hawful70/shop-identity-service/internal/identity/repository"
)

var (
	ErrEmailTaken      = errors.New("email is already registered")
	ErrInvalidLogin    = errors.New("invalid email or password")
	ErrPasswordTooWeak = errors.New("password must be at least 8 characters")
)

type Service interface {
	Register(ctx context.Context, email, username, password string) (User, error)
	Login(ctx context.Context, email, password string) (User, string, error)
}

type service struct {
	repo       repository.Repository
	jwtManager *JWTManager
}

func NewService(repo repository.Repository, jwtManager *JWTManager) Service {
	return &service{repo: repo, jwtManager: jwtManager}
}

func (s *service) Register(ctx context.Context, email, username, password string) (User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	username = strings.TrimSpace(username)

	if len(password) < 8 {
		return User{}, ErrPasswordTooWeak
	}

	_, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return User{}, ErrEmailTaken
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return User{}, err
	}

	hashed, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}

	user, err := NewUser(email, username, hashed)
	if err != nil {
		return User{}, err
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *service) Login(ctx context.Context, email, password string) (User, string, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return User{}, "", ErrInvalidLogin
		}
		return User{}, "", err
	}

	if !CheckPassword(user.Password, password) {
		return User{}, "", ErrInvalidLogin
	}

	token, err := s.jwtManager.GenerateToken(user)
	if err != nil {
		return User{}, "", err
	}

	return user, token, nil
}
