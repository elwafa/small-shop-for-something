package postgres

import (
	"context"
	"database/sql"
	"github.com/elwafa/billion-data/internal/entities"
)

type ItemRepo struct {
	DB *sql.DB
}

func NewPostgresItemRepository(db *sql.DB) *ItemRepo {
	return &ItemRepo{
		DB: db,
	}
}

func (r *ItemRepo) StoreItem(ctx context.Context, item *entities.Item) (*entities.Item, error) {
	// insert item to database and return the id
	err := r.DB.QueryRowContext(ctx, "INSERT INTO items(name, description, price, picture, status, receive, user_id) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		item.Name, item.Description, item.Price, item.Picture, item.Status, item.Receive, item.UserId).Scan(&item.ID)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (r *ItemRepo) GetItems(ctx context.Context, limit int) ([]entities.Item, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT * FROM items LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entities.Item
	for rows.Next() {
		var item entities.Item
		err = rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Receive, &item.UserId)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepo) GetItemsByUser(ctx context.Context, userId, limit int) ([]entities.Item, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT * FROM items WHERE user_id=$1 LIMIT $2", userId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entities.Item
	for rows.Next() {
		var item entities.Item
		err = rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Receive, &item.UserId)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepo) GetItem(ctx context.Context, itemID int) (entities.Item, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT * FROM items WHERE id=$1", itemID)
	var item entities.Item
	err := row.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Receive, &item.UserId)
	if err != nil {
		return entities.Item{}, err
	}
	return item, nil
}

func (r *ItemRepo) GetItemByUser(ctx context.Context, userId, itemID int) (entities.Item, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT * FROM items WHERE user_id=$1 AND id=$2", userId, itemID)
	var item entities.Item
	err := row.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Receive, &item.UserId)
	if err != nil {
		return entities.Item{}, err
	}
	return item, nil
}

func (r *ItemRepo) UpdateItem(ctx context.Context, itemID int, item entities.Item) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE items SET name=$1, description=$2, price=$3, picture=$4, status=$5, receive=$6 WHERE id=$7",
		item.Name, item.Description, item.Price, item.Picture, item.Status, item.Receive, itemID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ItemRepo) DeleteItem(ctx context.Context, itemID int) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM items WHERE id=$1", itemID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ItemRepo) GetPaginationItems(ctx context.Context, limit, page int, sort, name string) ([]entities.Item, int, error) {
	// check if name is empty
	query := "SELECT id, name,description,price,picture,status,receive,user_id FROM items"
	rows := &sql.Rows{}
	var err error
	var total int
	var totalCount *sql.Row
	if name != "" {
		query += " WHERE name like $1 ORDER BY created_at " + sort + " LIMIT $2 OFFSET $3"
		rows, err = r.DB.QueryContext(ctx, query, "%"+name+"%", limit, (page-1)*limit)
		totalCount = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM items WHERE name like $1 ORDER BY $2 LIMIT $3 OFFSET $4", "%"+name+"%", sort, limit, (page-1)*limit)
	} else {
		query += " ORDER BY created_at " + sort + " LIMIT $1 OFFSET $2"
		rows, err = r.DB.QueryContext(ctx, query, limit, (page-1)*limit)
		totalCount = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM items ORDER BY $1 LIMIT $2 OFFSET $3", sort, limit, (page-1)*limit)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var items []entities.Item
	for rows.Next() {
		var item entities.Item
		err = rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Receive, &item.UserId)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	err = totalCount.Scan(&total)
	return items, total, nil
}

func (r *ItemRepo) GetPaginationItemsByUser(ctx context.Context, userId, limit, page int) ([]entities.Item, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, name,description,price,picture,status,receive,user_id FROM items WHERE user_id=$1 LIMIT $2 OFFSET $3", userId, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entities.Item
	for rows.Next() {
		var item entities.Item
		err = rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Receive, &item.UserId)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepo) GetTotalItemsByUser(ctx context.Context, userId int) (int, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM items WHERE user_id=$1", userId)
	var total int
	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (r *ItemRepo) GetTotalItems(ctx context.Context) (int, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM items")
	var total int
	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
