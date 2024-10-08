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

func (r *OrderRepository) GetOrders(ctx context.Context, userId int) ([]*entities.Order, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, created_by, status FROM orders WHERE created_by = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := make([]*entities.Order, 0)
	for rows.Next() {
		order := &entities.Order{}
		err = rows.Scan(&order.ID, &order.UserID, &order.Status)
		if err != nil {
			return nil, err
		}
		// get order items
		orderItems, err := r.getOrderItemsAndItem(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.OrderItems = orderItems
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) getOrderItemsAndItem(ctx context.Context, orderId int) ([]*entities.OrderItem, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, item_id, status FROM order_items WHERE order_id = $1", orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orderItems := make([]*entities.OrderItem, 0)
	for rows.Next() {
		orderItem := &entities.OrderItem{}
		err = rows.Scan(&orderItem.ID, &orderItem.ItemID, &orderItem.Status)
		if err != nil {
			return nil, err
		}
		// get item
		item, err := r.getItem(ctx, orderItem.ItemID)
		if err != nil {
			return nil, err
		}
		orderItem.Item = item
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}

func (r *OrderRepository) getItem(ctx context.Context, itemId int) (*entities.Item, error) {
	item := &entities.Item{}
	err := r.db.QueryRowContext(ctx, "SELECT id, name, price,picture FROM items WHERE id = $1", itemId).
		Scan(&item.ID, &item.Name, &item.Price, &item.Picture)
	if err != nil {
		return nil, err
	}
	return item, nil
}
