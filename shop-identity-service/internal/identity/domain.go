package identity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserID string

type User struct {
	ID        UserID
	Email     string
	Username  string
	Password  string // hashed
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrInvalidUser   = errors.New("invalid user")
	ErrEmailRequired = errors.New("email is required")
	ErrPasswordShort = errors.New("password too short")
)

func NewUser(email, username, hashedPassword string) (User, error) {
	if email == "" {
		return User{}, ErrEmailRequired
	}
	if len(hashedPassword) == 0 {
		return User{}, ErrInvalidUser
	}

	now := time.Now().UTC()
	return User{
		ID:        UserID(uuid.NewString()),
		Email:     email,
		Username:  username,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
