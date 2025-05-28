package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type OrderItem struct {
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Price     int    `json:"price"`
	Quantity  int    `json:"quantity"`
	Subtotal  int    `json:"subtotal"`
}

type Order struct {
	ID        int         `json:"id"`
	UserID    string      `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Total     int         `json:"total"`
	Timestamp string      `json:"timestamp"`
}

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) PlaceOrder(userID string) (*Order, error) {
	// Get user's cart
	rows, err := r.DB.Query(`
		SELECT c.product_id, p.name, p.price, c.quantity
		FROM carts c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %v", err)
	}
	defer rows.Close()

	var items []OrderItem
	var total int
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ProductID, &item.Name, &item.Price, &item.Quantity); err != nil {
			return nil, err
		}
		item.Subtotal = item.Price * item.Quantity
		total += item.Subtotal
		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Insert order
	itemBytes, err := json.Marshal(items)
	if err != nil {
		return nil, fmt.Errorf("failed to encode order items")
	}

	var orderID int
	var timestamp string
	err = r.DB.QueryRow(`
		INSERT INTO orders (user_id, items, total)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, userID, itemBytes, total).Scan(&orderID, &timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to insert order: %v", err)
	}

	// Clear cart
	_, _ = r.DB.Exec("DELETE FROM carts WHERE user_id = $1", userID)

	return &Order{
		ID:        orderID,
		UserID:    userID,
		Items:     items,
		Total:     total,
		Timestamp: timestamp,
	}, nil
}

func (r *OrderRepository) GetOrderHistory(userID string) ([]Order, error) {
	rows, err := r.DB.Query(`
		SELECT id, items, total, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders")
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		var itemBytes []byte
		if err := rows.Scan(&o.ID, &itemBytes, &o.Total, &o.Timestamp); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(itemBytes, &o.Items); err != nil {
			return nil, fmt.Errorf("failed to decode order items")
		}
		orders = append(orders, o)
	}
	return orders, nil
}
