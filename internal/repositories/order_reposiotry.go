package repositories

import (
	"context"
	"github.com/elwafa/billion-data/internal/entities"
)

type OrderRepository interface {
	StoreOrder(ctx context.Context, order *entities.Order) (*entities.Order, error)
	GetOrders(ctx context.Context, userId int) ([]*entities.Order, error)
	UpdateOrderItemStatus(ctx context.Context, orderId, itemId int) error
	GetOrderItems(ctx context.Context, orderId int) ([]*entities.OrderItem, error)
	UpdateOrderStatus(ctx context.Context, orderId int) error
	GetItemsForSeller(ctx context.Context, userId int) ([]*entities.OrderItem, error)
}
