package controllers

import (
	"bego/middleware"
	"bego/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// POST /api/cart
func AddToCart(w http.ResponseWriter, r *http.Request) {
	var cartService = services.NewCartService()
	role := r.Context().Value(middleware.ContextRole).(string)
	userID := r.Context().Value(middleware.ContextUserID).(string)

	if role != "user" {
		http.Error(w, "Only users can access cart", http.StatusForbidden)
		return
	}

	var input struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Quantity <= 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	item, err := cartService.AddToCart(userID, input.ProductID, input.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(item)
}

// GET /api/cart
func ViewCart(w http.ResponseWriter, r *http.Request) {
	var cartService = services.NewCartService()
	role := r.Context().Value(middleware.ContextRole).(string)
	userID := r.Context().Value(middleware.ContextUserID).(string)

	if role != "user" {
		http.Error(w, "Only users can view cart", http.StatusForbidden)
		return
	}

	cart, err := cartService.GetCart(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cart)
}

// PUT /api/cart/{id}
func UpdateCartItem(w http.ResponseWriter, r *http.Request) {

	var cartService = services.NewCartService()
	role := r.Context().Value(middleware.ContextRole).(string)
	if role != "user" {
		http.Error(w, "Only users can update cart", http.StatusForbidden)
		return
	}

	cartID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Quantity int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Quantity <= 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	err = cartService.UpdateQuantity(cartID, input.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Cart updated"})
}

// DELETE /api/cart/{id}
func DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	var cartService = services.NewCartService()
	role := r.Context().Value(middleware.ContextRole).(string)
	if role != "user" {
		http.Error(w, "Only users can delete cart items", http.StatusForbidden)
		return
	}

	cartID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid cart ID", http.StatusBadRequest)
		return
	}

	err = cartService.RemoveFromCart(cartID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Item removed from cart"})
}
