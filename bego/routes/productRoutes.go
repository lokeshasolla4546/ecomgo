package routes

import (
	"bego/controllers"
	"bego/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterProductRoutes(router *mux.Router) {
	// Public: Get all products
	router.HandleFunc("/api/products", controllers.GetAllProducts).Methods("GET")

	// Admin-only: Add, Delete, Update
	router.Handle("/api/products", middleware.ProtectRoutes(http.HandlerFunc(controllers.AddProduct))).Methods("POST")
	router.Handle("/api/products/{id}", middleware.ProtectRoutes(http.HandlerFunc(controllers.DeleteProduct))).Methods("DELETE")
	router.Handle("/api/products/{id}", middleware.ProtectRoutes(http.HandlerFunc(controllers.UpdateProduct))).Methods("PUT")
}
