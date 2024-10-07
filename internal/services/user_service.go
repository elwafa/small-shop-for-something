package services

import (
	"context"
	"errors"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/repositories"
)

// Error messages
var (
	ErrUserAlreadyExist = errors.New("user already exist")
)

type UserService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) StoreUser(ctx context.Context, user *entities.User) error {
	checkUser, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err == nil && user.Email == checkUser.Email {
		return ErrUserAlreadyExist
	}
	checkPhone, err := s.repo.GetUserByPhone(ctx, user.Phone)
	if err == nil && user.Phone == checkPhone.Phone {
		return errors.New("phone number already exist")
	}
	err = s.repo.StoreUser(ctx, *user)
	if err != nil {
		return err
	}
	return nil
}
