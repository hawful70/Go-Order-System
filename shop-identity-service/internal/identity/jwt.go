package identity

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret    []byte
	issuer    string
	expiresIn time.Duration
}

type Claims struct {
	UserID   string `json:"uid"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret, issuer string, expiresIn time.Duration) *JWTManager {
	return &JWTManager{
		secret:    []byte(secret),
		issuer:    issuer,
		expiresIn: expiresIn,
	}
}

func (m *JWTManager) GenerateToken(u User) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		UserID:   string(u.ID),
		Email:    u.Email,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   string(u.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *JWTManager) VerifyToken(tokenStr string) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	}, jwt.WithIssuer(m.issuer))
	if err != nil {
		return Claims{}, err
	}
	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}
	return claims, nil
}
