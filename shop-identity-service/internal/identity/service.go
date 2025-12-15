package identity

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hawful70/shop-identity-service/internal/identity/repository"
)

var (
	ErrEmailTaken      = errors.New("email is already registered")
	ErrInvalidLogin    = errors.New("invalid email or password")
	ErrPasswordTooWeak = errors.New("password must be at least 8 characters")
	ErrInvalidToken    = errors.New("invalid token")
)

type Service interface {
	Register(ctx context.Context, email, username, password string) (User, error)
	Login(ctx context.Context, email, password string) (User, string, error)
	GetUserByID(ctx context.Context, id UserID) (User, error)
	ValidateToken(ctx context.Context, token string) (User, Claims, error)
}

type service struct {
	repo       repository.Repository
	jwtManager *JWTManager
	notifier   UserNotifier
}

func NewService(repo repository.Repository, jwtManager *JWTManager, notifier UserNotifier) Service {
	if notifier == nil {
		notifier = NoopNotifier()
	}
	return &service{repo: repo, jwtManager: jwtManager, notifier: notifier}
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

	if err := s.notifier.UserCreated(ctx, user); err != nil {
		fmt.Println("err", err)
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

func (s *service) GetUserByID(ctx context.Context, id UserID) (User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *service) ValidateToken(ctx context.Context, token string) (User, Claims, error) {
	claims, err := s.jwtManager.VerifyToken(token)
	if err != nil {
		return User{}, Claims{}, ErrInvalidToken
	}

	user, err := s.repo.GetUserByID(ctx, UserID(claims.UserID))
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return User{}, claims, ErrInvalidToken
		}
		return User{}, claims, err
	}

	return user, claims, nil
}
