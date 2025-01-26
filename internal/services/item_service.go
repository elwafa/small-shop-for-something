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

func (s *ItemService) GetItemsForSeller(ctx context.Context, limit, page, userId int) ([]entities.Item, int, error) {
	items, err := s.repo.GetPaginationItemsByUser(ctx, userId, limit, page)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.GetTotalItemsByUser(ctx, userId)
	return items, total, nil
}

func (s *ItemService) GetItemsForCustomer(ctx context.Context, limit, page int, sort, name, colour, category string, price float64) ([]entities.Item, int, error) {
	items, total, err := s.repo.GetPaginationItems(ctx, limit, page, sort, colour, category, name, price)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (s *ItemService) GetItemForCustomer(ctx context.Context, id int) (entities.Item, error) {
	item, err := s.repo.GetItem(ctx, id)
	if err != nil {
		return entities.Item{}, err
	}
	return *item, nil
}

func (s *ItemService) DeleteItem(ctx context.Context, id, userId int) error {
	// Validate if the item belongs to the user
	item, err := s.repo.GetItemByUser(ctx, userId, id)
	if err != nil {
		return err
	}
	if item.ID == 0 {
		return repositories.ErrRecodeNotFound
	}
	return s.repo.DeleteItem(ctx, id)
}
