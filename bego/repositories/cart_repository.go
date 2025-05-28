package repositories

import (
	"database/sql"
	"fmt"
)

type CartItem struct {
	ID        int
	UserID    string
	ProductID int
	Quantity  int
}

type EnrichedCartItem struct {
	CartID    int
	ProductID int
	Name      string
	Price     int
	Quantity  int
	Image     string
}

type CartRepository struct {
	DB *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{DB: db}
}

func (r *CartRepository) Add(userID string, productID, quantity int) (*CartItem, error) {
	var id int
	var existingQty, productQty int

	// Check stock
	err := r.DB.QueryRow("SELECT quantity FROM products WHERE id = $1", productID).Scan(&productQty)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}
	if productQty < quantity {
		return nil, fmt.Errorf("not enough stock available")
	}

	// Check existing cart item
	err = r.DB.QueryRow("SELECT id, quantity FROM carts WHERE user_id = $1 AND product_id = $2", userID, productID).Scan(&id, &existingQty)
	if err == nil {
		newQty := existingQty + quantity
		if productQty < (newQty - existingQty) {
			return nil, fmt.Errorf("not enough stock to add more")
		}
		_, err = r.DB.Exec("UPDATE carts SET quantity = $1 WHERE id = $2", newQty, id)
		if err != nil {
			return nil, fmt.Errorf("failed to update cart")
		}
		_, _ = r.DB.Exec("UPDATE products SET quantity = quantity - $1 WHERE id = $2", quantity, productID)
		return &CartItem{ID: id, UserID: userID, ProductID: productID, Quantity: newQty}, nil
	}

	// Insert new cart item
	err = r.DB.QueryRow("INSERT INTO carts (user_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING id", userID, productID, quantity).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to add to cart")
	}
	_, _ = r.DB.Exec("UPDATE products SET quantity = quantity - $1 WHERE id = $2", quantity, productID)
	return &CartItem{ID: id, UserID: userID, ProductID: productID, Quantity: quantity}, nil
}

func (r *CartRepository) Get(userID string) ([]EnrichedCartItem, error) {
	rows, err := r.DB.Query(`
		SELECT c.id, c.product_id, p.name, p.price, c.quantity, p.image
		FROM carts c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart")
	}
	defer rows.Close()

	var cart []EnrichedCartItem
	for rows.Next() {
		var item EnrichedCartItem
		if err := rows.Scan(&item.CartID, &item.ProductID, &item.Name, &item.Price, &item.Quantity, &item.Image); err != nil {
			return nil, err
		}
		cart = append(cart, item)
	}
	return cart, nil
}

func (r *CartRepository) UpdateQuantity(cartID, newQty int) error {
	var oldQty, productID, stockQty int
	err := r.DB.QueryRow("SELECT quantity, product_id FROM carts WHERE id = $1", cartID).Scan(&oldQty, &productID)
	if err != nil {
		return fmt.Errorf("cart item not found")
	}
	delta := newQty - oldQty
	if delta > 0 {
		err = r.DB.QueryRow("SELECT quantity FROM products WHERE id = $1", productID).Scan(&stockQty)
		if err != nil || stockQty < delta {
			return fmt.Errorf("not enough stock to increase quantity")
		}
	}

	_, err = r.DB.Exec("UPDATE carts SET quantity = $1 WHERE id = $2", newQty, cartID)
	if err != nil {
		return fmt.Errorf("failed to update cart")
	}
	_, _ = r.DB.Exec("UPDATE products SET quantity = quantity - $1 WHERE id = $2", delta, productID)
	return nil
}

func (r *CartRepository) Remove(cartID int) error {
	var productID, qty int
	err := r.DB.QueryRow("SELECT product_id, quantity FROM carts WHERE id = $1", cartID).Scan(&productID, &qty)
	if err != nil {
		return fmt.Errorf("cart item not found")
	}
	_, err = r.DB.Exec("DELETE FROM carts WHERE id = $1", cartID)
	if err != nil {
		return fmt.Errorf("failed to delete cart item")
	}
	_, _ = r.DB.Exec("UPDATE products SET quantity = quantity + $1 WHERE id = $2", qty, productID)
	return nil
}
