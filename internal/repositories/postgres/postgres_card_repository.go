package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/repositories"
)

type CardRepository struct {
	db *sql.DB
}

func NewPostgresCardRepository(db *sql.DB) *CardRepository {
	return &CardRepository{
		db: db,
	}
}

func (r *CardRepository) CheckPendingCardForUser(ctx context.Context, userId int) (*entities.Card, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id,user_id,status FROM card WHERE user_id = $1 AND status = 'pending'", userId)
	card := &entities.Card{}
	err := row.Scan(&card.ID, &card.UserID, &card.Status)
	if err != nil {
		// if no pending card found, error will be sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrRecodeNotFound
		}
		return nil, err
	}
	return card, nil
}

func (r *CardRepository) StoreCard(ctx context.Context, card *entities.Card) (*entities.Card, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	err = tx.QueryRowContext(ctx, "INSERT INTO card(user_id,status) VALUES($1,$2) RETURNING id", card.UserID, card.Status).
		Scan(&card.ID)
	if err != nil {
		errRollback := tx.Rollback()
		fmt.Println(err, "rollback")
		if errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (r *CardRepository) AddToCard(ctx context.Context, cardItem *entities.CardItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO card_item(card_id,item_id) VALUES($1,$2)", cardItem.CardID, cardItem.ItemID)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r *CardRepository) CheckItemExistsInCard(ctx context.Context, cardId, itemId int) (*entities.CardItem, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id,card_id,item_id FROM card_item WHERE card_id = $1 AND item_id = $2", cardId, itemId)
	cardItem := &entities.CardItem{}
	err := row.Scan(&cardItem.ID, &cardItem.CardID, &cardItem.ItemID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrRecodeNotFound
		}
		return nil, err
	}
	return cardItem, nil
}

func (r *CardRepository) GetCard() error {
	return nil
}

func (r *CardRepository) RemoveFromCard() error {
	return nil
}
