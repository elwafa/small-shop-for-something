package services

import (
	"context"
	"errors"
	"github.com/elwafa/billion-data/internal/auth"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/repositories"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidCredential = errors.New("invalid credential")
)

type AuthService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (entities.User, error, string) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repositories.ErrRecodeNotFound) {
			return user, ErrInvalidCredential, ""
		}
		return user, err, ""
	}
	checkPass := entities.CheckPasswordHash(password, user.Password)
	if !checkPass {
		return user, ErrInvalidCredential, ""
	}
	token, err := auth.Login(user)
	if err != nil {
		return user, err, ""
	}
	return user, nil, token
}

// admin login
func (s *AuthService) AdminLogin(ctx *gin.Context, email, password string) (entities.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repositories.ErrRecodeNotFound) {
			return user, ErrInvalidCredential
		}
		return user, err
	}
	checkPass := entities.CheckPasswordHash(password, user.Password)
	if !checkPass {
		return user, ErrInvalidCredential
	}
	if user.Type != "admin" {
		return user, ErrInvalidCredential
	}
	return user, auth.AdminLoginWeb(ctx, user)
}
