package repositories

import (
	"context"
	"github.com/elwafa/billion-data/internal/entities"
)

type ItemRepository interface {
	StoreItem(ctx context.Context, item *entities.Item) (*entities.Item, error)
	GetItems(ctx context.Context, limit int) ([]entities.Item, error)
	GetItemsByUser(ctx context.Context, userId, limit int) ([]entities.Item, error)
	GetItem(ctx context.Context, itemID int) (entities.Item, error)
	GetItemByUser(ctx context.Context, userId, itemID int) (entities.Item, error)
	UpdateItem(ctx context.Context, itemID int, item entities.Item) error
	DeleteItem(ctx context.Context, itemID int) error
	GetPaginationItems(ctx context.Context, limit, page int) ([]entities.Item, error)
	GetPaginationItemsByUser(ctx context.Context, userId, limit, page int) ([]entities.Item, error)
	GetTotalItemsByUser(ctx context.Context, userId int) (int, error)
}
