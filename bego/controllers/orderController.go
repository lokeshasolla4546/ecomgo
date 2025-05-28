package controllers

import (
	"bego/middleware"
	"bego/services"
	"encoding/json"
	"net/http"
)

// POST /api/orders
func CreateOrder(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.ContextRole).(string)
	userID := r.Context().Value(middleware.ContextUserID).(string)
	var orderService = services.NewOrderService()
	if role != "user" {
		http.Error(w, "Only users can place orders", http.StatusForbidden)
		return
	}

	order, err := orderService.PlaceOrder(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Payment Successful",
		"order":   order,
	})
}

// GET /api/orders
func GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value(middleware.ContextRole).(string)
	userID := r.Context().Value(middleware.ContextUserID).(string)
	var orderService = services.NewOrderService()
	if role != "user" {
		http.Error(w, "Only users can view orders", http.StatusForbidden)
		return
	}

	orders, err := orderService.GetOrderHistory(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}
