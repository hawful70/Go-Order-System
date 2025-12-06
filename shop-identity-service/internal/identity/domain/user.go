package domain

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
	Password  string
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

type UserModel struct {
	ID        string `gorm:"primaryKey;type:text"`
	Email     string `gorm:"uniqueIndex;type:text"`
	Username  string `gorm:"type:text"`
	Password  string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserModel) TableName() string {
	return "users"
}

func ToUserModel(u User) UserModel {
	return UserModel{
		ID:        string(u.ID),
		Email:     u.Email,
		Username:  u.Username,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (m UserModel) ToDomain() User {
	return User{
		ID:        UserID(m.ID),
		Email:     m.Email,
		Username:  m.Username,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
