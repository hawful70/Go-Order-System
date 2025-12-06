package identity

import (
	"context"
	"errors"
	"strings"
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
	repo       Repository
	jwtManager *JWTManager
}

func NewService(repo Repository, jwtManager *JWTManager) Service {
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
		// user exists
		return User{}, ErrEmailTaken
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
		return User{}, "", ErrInvalidLogin
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
