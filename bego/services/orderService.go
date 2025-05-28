package services

import (
	"bego/config"
	"bego/repositories"
)

type OrderService struct {
	repo *repositories.OrderRepository
}

func NewOrderService() *OrderService {
	return &OrderService{
		repo: repositories.NewOrderRepository(config.DB),
	}
}

func (s *OrderService) PlaceOrder(userID string) (*repositories.Order, error) {
	return s.repo.PlaceOrder(userID)
}

func (s *OrderService) GetOrderHistory(userID string) ([]repositories.Order, error) {
	return s.repo.GetOrderHistory(userID)
}
