package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/halo-dev/halo-go/internal/model"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	Register(ctx context.Context, user *model.User) (string, error)
	RefreshToken(_ context.Context, tokenString string) (string, error)
	GenerateToken(username string) (string, error)
}

type authService struct {
	userSvc UserService
	secret  string
	expire  int
}

func NewAuthService(userSvc UserService, secret string, expireHours int) AuthService {
	return &authService{
		userSvc: userSvc,
		secret:  secret,
		expire:  expireHours,
	}
}

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userSvc.GetByName(ctx, username)
	if err != nil {
		return "", fmt.Errorf("用户名或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Spec.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("用户名或密码错误")
	}
	return s.GenerateToken(user.Spec.UserName)
}

func (s *authService) Register(ctx context.Context, user *model.User) (string, error) {
	if _, err := s.userSvc.Create(ctx, user); err != nil {
		return "", fmt.Errorf("注册失败: %w", err)
	}
	return s.GenerateToken(user.Spec.UserName)
}

func (s *authService) RefreshToken(_ context.Context, tokenString string) (string, error) {
	username, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("令牌无效，无法刷新")
	}
	return s.GenerateToken(username)
}

func (s *authService) GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Duration(s.expire) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *authService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.secret), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("令牌无效或已过期")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("无法解析令牌声明")
	}
	username, _ := claims["sub"].(string)
	return username, nil
}
