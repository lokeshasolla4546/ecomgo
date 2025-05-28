package repositories

import (
	"database/sql"
	"fmt"
)

type Product struct {
	ID       int
	Name     string
	Price    int
	Quantity int
	Image    string
}

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) Add(name string, price, quantity int, image string) (*Product, error) {
	var id int
	err := r.DB.QueryRow(
		"INSERT INTO products (name, price, quantity, image) VALUES ($1, $2, $3, $4) RETURNING id",
		name, price, quantity, image,
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert product: %v", err)
	}
	return &Product{ID: id, Name: name, Price: price, Quantity: quantity, Image: image}, nil
}

func (r *ProductRepository) GetAll() ([]Product, error) {
	rows, err := r.DB.Query("SELECT id, name, price, quantity, image FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Image); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}

func (r *ProductRepository) Update(id int, name string, price int, quantity int, image string) (*Product, error) {
	_, err := r.DB.Exec(
		"UPDATE products SET name = $1, price = $2, quantity = $3, image = $4 WHERE id = $5",
		name, price, quantity, image, id,
	)
	if err != nil {
		return nil, err
	}
	return &Product{ID: id, Name: name, Price: price, Quantity: quantity, Image: image}, nil
}
