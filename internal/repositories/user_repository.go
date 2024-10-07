package repositories

import (
	"context"
	"github.com/elwafa/billion-data/internal/entities"
)

type UserRepository interface {
	StoreUser(ctx context.Context, user entities.User) error
	GetUser(ctx context.Context, userID int) (entities.User, error)
	GetUsers(ctx context.Context, limit, offset int, isActiveOnly bool) ([]entities.User, error)
	UpdateUser(ctx context.Context, userID int, user entities.User) error
	DeleteUser(ctx context.Context, userID int) error
	// Get User by email
	GetUserByEmail(ctx context.Context, email string) (entities.User, error)
	// Get User by phone
	GetUserByPhone(ctx context.Context, phone string) (entities.User, error)
}
