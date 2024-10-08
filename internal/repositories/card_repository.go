package repositories

import (
	"context"
	"github.com/elwafa/billion-data/internal/entities"
)

type CardRepository interface {
	CheckPendingCardForUser(ctx context.Context, userId int) (*entities.Card, error)
	StoreCard(ctx context.Context, card *entities.Card) (*entities.Card, error)
	AddToCard(ctx context.Context, cardItem *entities.CardItem) error
	CheckItemExistsInCard(ctx context.Context, cardId, itemId int) (*entities.CardItem, error)
	GetCard() error
	RemoveFromCard() error
}
