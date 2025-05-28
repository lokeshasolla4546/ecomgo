package services

import (
	"bego/config"
	"bego/repositories"
)

type ProductService struct {
	repo *repositories.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{
		repo: repositories.NewProductRepository(config.DB),
	}
}

func (s *ProductService) AddProduct(name string, price int, quantity int, image string) (*repositories.Product, error) {
	return s.repo.Add(name, price, quantity, image)
}

func (s *ProductService) GetAllProducts() ([]repositories.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductService) UpdateProduct(id int, name string, price int, quantity int, image string) (*repositories.Product, error) {
	return s.repo.Update(id, name, price, quantity, image)
}
