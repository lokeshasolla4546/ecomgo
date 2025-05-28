package routes

import (
	"bego/controllers"
	"bego/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterCartRoutes(router *mux.Router) {
	// User-only protected routes
	router.Handle("/api/cart", middleware.ProtectRoutes(http.HandlerFunc(controllers.AddToCart))).Methods("POST")
	router.Handle("/api/cart", middleware.ProtectRoutes(http.HandlerFunc(controllers.ViewCart))).Methods("GET")
	router.Handle("/api/cart/{id}", middleware.ProtectRoutes(http.HandlerFunc(controllers.UpdateCartItem))).Methods("PUT")
	router.Handle("/api/cart/{id}", middleware.ProtectRoutes(http.HandlerFunc(controllers.DeleteCartItem))).Methods("DELETE")
}
