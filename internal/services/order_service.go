package services

import (
	"context"
	"errors"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/repositories"
	"github.com/gin-gonic/gin"
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

func (s *OrderService) GetOrders(ctx context.Context, userId int) ([]*entities.Order, error) {
	orders, err := s.repo.GetOrders(ctx, userId)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// Update Item Status In Order to be Done
func (s *OrderService) UpdateItemStatus(ctx *gin.Context, orderId, itemId int) error {
	err := s.repo.UpdateOrderItemStatus(ctx, orderId, itemId)
	if err != nil {
		return err
	}
	return nil
}

// Update Order Status to be Done

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderId int) error {
	// first check if all items in order are done
	orderItems, err := s.repo.GetOrderItems(ctx, orderId)
	if err != nil {
		return err
	}
	for _, orderItem := range orderItems {
		if orderItem.Status != "delivered" {
			return nil
		}
	}
	// update order status to be delivered
	return s.repo.UpdateOrderStatus(ctx, orderId)
}

// get order items for seller which are not delivered yet and assign it to him
func (s *OrderService) GetOrderItemsForSeller(ctx context.Context, userId int) ([]*entities.OrderItem, error) {
	// get all orders for user
	orderItems, err := s.repo.GetItemsForSeller(ctx, userId)
	if err != nil {
		return nil, err
	}
	//get all order items which are not delivered yet
	return orderItems, nil
}
