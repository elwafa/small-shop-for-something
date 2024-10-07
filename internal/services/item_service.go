package services

import (
	"context"
	"github.com/elwafa/billion-data/internal/entities"
	"github.com/elwafa/billion-data/internal/repositories"
)

type ItemService struct {
	repo repositories.ItemRepository
}

func NewItemService(repo repositories.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) StoreItem(ctx context.Context, item *entities.Item) (*entities.Item, error) {
	itemDb, err := s.repo.StoreItem(ctx, item)
	if err != nil {
		return &entities.Item{}, err
	}
	return itemDb, nil
}
