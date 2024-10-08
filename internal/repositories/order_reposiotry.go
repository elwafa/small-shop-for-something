package repositories

import (
	"context"
	"github.com/elwafa/billion-data/internal/entities"
)

type OrderRepository interface {
	StoreOrder(ctx context.Context, order *entities.Order) (*entities.Order, error)
	GetOrders(ctx context.Context, userId int) ([]*entities.Order, error)
}
