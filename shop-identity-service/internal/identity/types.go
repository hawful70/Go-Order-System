package identity

import domain "github.com/hawful70/shop-identity-service/internal/identity/domain"

type User = domain.User
type UserID = domain.UserID

var (
	ErrInvalidUser   = domain.ErrInvalidUser
	ErrEmailRequired = domain.ErrEmailRequired
	ErrPasswordShort = domain.ErrPasswordShort
)

func NewUser(email, username, hashedPassword string) (User, error) {
	return domain.NewUser(email, username, hashedPassword)
}
