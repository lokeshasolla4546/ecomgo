package controllers

import (
	"bego/middleware"
	"bego/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var productService = services.NewProductService()
	role := r.Context().Value(middleware.ContextRole).(string)
	if role != "admin" {
		http.Error(w, "Access denied: admin only", http.StatusForbidden)
		return
	}

	var input struct {
		Name     string `json:"name"`
		Price    int    `json:"price"`
		Quantity int    `json:"quantity"`
		Image    string `json:"image"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Name == "" || input.Price <= 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	product, err := productService.AddProduct(input.Name, input.Price, input.Quantity, input.Image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(product)
}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	var productService = services.NewProductService()
	products, err := productService.GetAllProducts()
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var productService = services.NewProductService()
	role := r.Context().Value(middleware.ContextRole).(string)
	if role != "admin" {
		http.Error(w, "Access denied: admin only", http.StatusForbidden)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := productService.DeleteProduct(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted"})
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var productService = services.NewProductService()
	role := r.Context().Value(middleware.ContextRole).(string)
	if role != "admin" {
		http.Error(w, "Access denied: admin only", http.StatusForbidden)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Name     string `json:"name"`
		Price    int    `json:"price"`
		Quantity int    `json:"quantity"`
		Image    string `json:"image"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Name == "" || input.Price <= 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	product, err := productService.UpdateProduct(id, input.Name, input.Price, input.Quantity, input.Image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}
