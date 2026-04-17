package service

import (
	"errors"
	"time"

	"halo/internal/config"
	"halo/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	users  *UserService
	secret []byte
}

type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(users *UserService, cfg config.Config) *AuthService {
	return &AuthService{users: users, secret: []byte(cfg.JWTSecret)}
}

func (s *AuthService) Login(username, password string) (string, *model.User, error) {
	user, err := s.users.Authenticate(username, password)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", nil, err
	}

	return signed, user, nil
}

func (s *AuthService) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
