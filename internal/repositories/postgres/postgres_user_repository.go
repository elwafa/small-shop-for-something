package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/repositories"
)

type UserRepository struct {
	DB *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) StoreUser(ctx context.Context, user entities.User) error {
	err := r.DB.QueryRowContext(ctx,
		"INSERT INTO users (name, email, password,phone , type, is_active) VALUES ($1, $2, $3, $4, $5, $6)",
		user.Name, user.Email, user.Password, user.Phone, user.Type, user.IsActive).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, userID int) (entities.User, error) {
	var user entities.User
	err := r.DB.QueryRowContext(ctx,
		"SELECT name, email, password, is_active FROM users WHERE id = $1",
		userID).Scan(&user.Name, &user.Email, &user.Password, &user.IsActive)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context, limit, offset int, isActiveOnly bool) ([]entities.User, error) {
	// handle with pagination
	var users []entities.User
	rows, err := r.DB.QueryContext(ctx,
		"SELECT name, email, password, is_active FROM users WHERE is_active = $1 LIMIT $2 OFFSET $3",
		isActiveOnly, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user entities.User
		err := rows.Scan(&user.Name, &user.Email, &user.Password, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userID int, user entities.User) error {
	_, err := r.DB.ExecContext(ctx,
		"UPDATE users SET name = $1, email = $2, password = $3, is_active = $4 WHERE id = $5",
		user.Name, user.Email, user.Password, user.IsActive, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID int) error {
	_, err := r.DB.ExecContext(ctx,
		"DELETE FROM users WHERE id = $1",
		userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (entities.User, error) {
	var user entities.User
	err := r.DB.QueryRowContext(ctx,
		"SELECT id, phone,name, email,type , password , is_active FROM users WHERE email = $1",
		email).Scan(&user.Id, &user.Phone, &user.Name, &user.Email, &user.Type, &user.Password, &user.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, repositories.ErrRecodeNotFound
		}
		return entities.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByPhone(ctx context.Context, phone string) (entities.User, error) {
	var user entities.User
	err := r.DB.QueryRowContext(ctx,
		"SELECT id, phone,name, email,type , password , is_active FROM users WHERE phone = $1",
		phone).Scan(&user.Id, &user.Phone, &user.Name, &user.Email, &user.Type, &user.Password, &user.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, repositories.ErrRecodeNotFound
		}
		return entities.User{}, err
	}
	return user, nil
}
