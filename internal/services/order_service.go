package services

import (
	"context"
	"errors"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/repositories"
)

type OrderService struct {
	repo     repositories.OrderRepository
	cardRepo repositories.CardRepository
}

func NewOrderService(repo repositories.OrderRepository, cardRepo repositories.CardRepository) *OrderService {
	return &OrderService{repo: repo, cardRepo: cardRepo}
}

func (s *OrderService) StoreOrder(ctx context.Context, userId int) (*entities.Order, int, error) {
	card, err := s.cardRepo.CheckPendingCardForUser(ctx, userId)
	if err != nil {
		return nil, 0, errors.New("no pending card found")
	}
	// get card items
	cardAndItems, err := s.cardRepo.GetCard(ctx, userId)
	if err != nil {
		return nil, 0, err
	}
	order := &entities.Order{
		UserID: cardAndItems.UserID,
		Status: "pending",
	}
	for _, cardItem := range cardAndItems.Items {
		orderItem := &entities.OrderItem{
			ItemID: cardItem.ItemID,
			Status: "pending",
		}
		order.OrderItems = append(order.OrderItems, orderItem)
	}
	order, err = s.repo.StoreOrder(ctx, order)
	if err != nil {
		return nil, 0, err
	}
	return order, card.ID, nil
}
