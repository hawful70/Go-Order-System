package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/hawful70/shop-identity-service/internal/identity/domain"
)

var ErrUserNotFound = errors.New("user not found")

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateUser(ctx context.Context, u domain.User) error {
	model := domain.ToUserModel(u)
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *postgresRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var model domain.UserModel
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetUserByID(ctx context.Context, id domain.UserID) (domain.User, error) {
	var model domain.UserModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, err
	}
	return model.ToDomain(), nil
}
