package postgres

import (
	"context"
	"database/sql"
	"github.com/elwafa/billion-data/internal/entities"
)

type OrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) StoreOrder(ctx context.Context, order *entities.Order) (*entities.Order, error) {
	// insert order to database and return the id
	err := r.db.QueryRowContext(ctx, "INSERT INTO orders(created_by, status) VALUES($1, $2) RETURNING id",
		order.UserID, order.Status).Scan(&order.ID)
	if err != nil {
		return nil, err
	}
	// insert order items to database
	for _, orderItem := range order.OrderItems {
		err = r.db.QueryRowContext(ctx, "INSERT INTO order_items(order_id, item_id, status) VALUES($1, $2, $3) RETURNING id",
			order.ID, orderItem.ItemID, orderItem.Status).Scan(&orderItem.ID)
		if err != nil {
			return nil, err
		}
	}
	return order, nil
}
