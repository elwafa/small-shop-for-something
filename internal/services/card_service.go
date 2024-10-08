package services

import (
	"context"
	"errors"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/repositories"
)

type CardService struct {
	repo repositories.CardRepository
}

func NewCardService(repo repositories.CardRepository) *CardService {
	return &CardService{
		repo: repo,
	}
}

func (s *CardService) AddToCard(ctx context.Context, itemId, userId int) error {
	// check if users has pending card
	card, err := s.repo.CheckPendingCardForUser(ctx, userId)
	if errors.Is(err, repositories.ErrRecodeNotFound) {
		// if no pending card found, create a new card
		card, err = s.createCardForUser(ctx, userId)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}
	}
	cardItem := &entities.CardItem{
		CardID: card.ID,
		ItemID: itemId,
	}
	// check if item already in card
	_, err = s.repo.CheckItemExistsInCard(ctx, card.ID, itemId)
	if errors.Is(err, repositories.ErrRecodeNotFound) {
		// if item not in card, add it
		return s.repo.AddToCard(ctx, cardItem)
	} else {
		if err != nil {
			return err
		}
		return errors.New("item already in card")
	}
}

func (s *CardService) createCardForUser(ctx context.Context, userId int) (*entities.Card, error) {
	card := &entities.Card{
		UserID: userId,
		Status: "pending",
	}
	card, err := s.repo.StoreCard(ctx, card)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (s *CardService) GetCard(ctx context.Context, userId int) (*entities.Card, error) {
	return s.repo.GetCard(ctx, userId)
}

func (s *CardService) CheckPendingCardForUser(ctx context.Context, userId int) (*entities.Card, error) {
	card, err := s.repo.CheckPendingCardForUser(ctx, userId)
	return card, err
}

func (s *CardService) UpdateCardToBeFinished(ctx context.Context, cardId int) error {
	return s.repo.UpdateCard(ctx, cardId)
}
