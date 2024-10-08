package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/repositories"
	"log"
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

func (r *CardRepository) GetCard(ctx context.Context, userId int) (*entities.Card, error) {
	// get card by user id and card Items
	card := &entities.Card{}
	row := r.db.QueryRowContext(ctx, "SELECT id,user_id,status FROM card WHERE user_id = $1 AND status = 'pending'", userId)
	err := row.Scan(&card.ID, &card.UserID, &card.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("no pending card found", userId)
			return nil, repositories.ErrRecodeNotFound
		}
		return nil, err
	}
	// get card items and items
	rows, err := r.db.QueryContext(ctx, "SELECT ci.id,ci.card_id,ci.item_id,i.name,i.description,i.price,i.picture,i.status,i.receive,i.user_id FROM card_item ci JOIN items i ON ci.item_id = i.id WHERE ci.card_id = $1", card.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cardItems []entities.CardItem
	for rows.Next() {
		var cardItem entities.CardItem
		var item entities.Item
		err = rows.Scan(&cardItem.ID, &cardItem.CardID, &cardItem.ItemID, &item.Name, &item.Description, &item.Price, &item.Picture, &item.Status, &item.Receive, &item.UserId)
		if err != nil {
			return nil, err
		}
		cardItem.Item = &item
		cardItems = append(cardItems, cardItem)
	}
	card.Items = cardItems
	return card, nil
}

func (r *CardRepository) RemoveFromCard() error {
	return nil
}

func (r *CardRepository) UpdateCard(ctx context.Context, id int) error {
	_, err := r.db.Exec("UPDATE card SET status = 'finished' WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
