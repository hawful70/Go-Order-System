package repository

import (
	"context"

	"github.com/hawful70/shop-identity-service/internal/identity/domain"
)

type Repository interface {
	CreateUser(ctx context.Context, u domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}
