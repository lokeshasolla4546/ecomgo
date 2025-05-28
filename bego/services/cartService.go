package services

import (
	"bego/config"
	"bego/repositories"
)

type CartService struct {
	repo *repositories.CartRepository
}

func NewCartService() *CartService {
	return &CartService{
		repo: repositories.NewCartRepository(config.DB),
	}
}

func (s *CartService) AddToCart(userID string, productID int, quantity int) (*repositories.CartItem, error) {
	return s.repo.Add(userID, productID, quantity)
}

func (s *CartService) GetCart(userID string) ([]repositories.EnrichedCartItem, error) {
	return s.repo.Get(userID)
}

func (s *CartService) UpdateQuantity(cartID int, newQty int) error {
	return s.repo.UpdateQuantity(cartID, newQty)
}

func (s *CartService) RemoveFromCart(cartID int) error {
	return s.repo.Remove(cartID)
}
