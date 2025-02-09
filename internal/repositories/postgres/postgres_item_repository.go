package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/repositories"
	"strings"
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
	err := r.DB.QueryRowContext(ctx, "INSERT INTO items(name, description, price, picture, status, colour, category, receive, user_id) VALUES($1, $2, $3, $4, $5, $6, $7,$8,$9) RETURNING id",
		item.Name, item.Description, item.Price, item.Picture, item.Status, item.Color, item.Category, item.Receive, item.UserId).Scan(&item.ID)
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
		err = rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Color, &item.Category, &item.Receive, &item.UserId)
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
		err = rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Color, &item.Category, &item.Receive, &item.UserId)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepo) GetItem(ctx context.Context, itemID int) (*entities.Item, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT id,name,description,price,picture,status,colour,category,receive,user_id FROM items WHERE id=$1", itemID)
	var item entities.Item
	err := row.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Color, &item.Category, &item.Receive, &item.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &entities.Item{}, repositories.ErrRecodeNotFound
		}
		return &entities.Item{}, err
	}
	return &item, nil
}

func (r *ItemRepo) GetItemByUser(ctx context.Context, userId, itemID int) (entities.Item, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT * FROM items WHERE user_id=$1 AND id=$2", userId, itemID)
	var item entities.Item
	err := row.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Color, &item.Category, &item.Receive, &item.UserId)
	if err != nil {
		return entities.Item{}, err
	}
	return item, nil
}

func (r *ItemRepo) UpdateItem(ctx context.Context, itemID int, item entities.Item) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE items SET name=$1, description=$2, price=$3, picture=$4, status=$5,colour=$6,category=$7, receive=$6 WHERE id=$7",
		item.Name, item.Description, item.Price, item.Picture, item.Status, item.Color, item.Category, item.Receive, itemID)
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

func (r *ItemRepo) GetPaginationItems(
	ctx context.Context,
	limit, page int,
	sort, name, colour, category string,
	price float64,
) ([]entities.Item, int, error) {

	// Base SELECT and COUNT queries:
	baseQuery := `
        SELECT id, name, description, price, picture, status, colour, category, receive, user_id
        FROM items
    `
	countQuery := `
        SELECT COUNT(*)
        FROM items
    `

	// We'll build conditions dynamically:
	var conditions []string
	var params []interface{}
	paramIndex := 1 // tracks placeholder index ($1, $2, etc.)

	// If name is not empty, filter by name (case-insensitive LIKE, for example):
	if name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", paramIndex))
		params = append(params, "%"+name+"%")
		paramIndex++
	}

	// If colour is not empty, filter by colour:
	if colour != "" {
		conditions = append(conditions, fmt.Sprintf("colour = $%d", paramIndex))
		params = append(params, colour)
		paramIndex++
	}

	// If category is not empty, filter by category:
	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", paramIndex))
		params = append(params, category)
		paramIndex++
	}

	// If price is nonzero, filter by exact price (or adjust logic if you need >=, <=, etc.):
	if price != 0 {
		conditions = append(conditions, fmt.Sprintf("price = $%d", paramIndex))
		params = append(params, price)
		paramIndex++
	}

	// If we have conditions, append them to our queries:
	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Sorting by the 'name' column, default to ASC or DESC
	// NOTE: Make sure `sort` is either "ASC" or "DESC" to avoid SQL injection
	baseQuery += fmt.Sprintf(" ORDER BY name %s", sort)

	// Add pagination (limit and offset)
	// paramIndex is currently at the next available placeholder
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	params = append(params, limit, (page-1)*limit)
	paramIndex += 2

	// 1) Query the items with filters + order + limit/offset
	rows, err := r.DB.QueryContext(ctx, baseQuery, params...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// 2) Query the total count (using the same conditions, but without limit/offset)
	//    We only pass the parameters related to the WHERE clause (not the last two for limit/offset).
	countRow := r.DB.QueryRowContext(ctx, countQuery, params[:len(params)-2]...)

	// Scan the count
	var total int
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, err
	}

	// 3) Read the items
	var items []entities.Item
	for rows.Next() {
		var item entities.Item
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Price,
			&item.Picture,
			&item.Status,
			&item.Color,
			&item.Category,
			&item.Receive,
			&item.UserId,
		)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (r *ItemRepo) GetPaginationItemsByUser(ctx context.Context, userId, limit, page int) ([]entities.Item, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, name,description,price,picture,status,colour,category,receive,user_id FROM items WHERE user_id=$1 LIMIT $2 OFFSET $3", userId, limit, (page-1)*limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []entities.Item
	for rows.Next() {
		var item entities.Item
		err = rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Color, &item.Category, &item.Receive, &item.UserId)
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
