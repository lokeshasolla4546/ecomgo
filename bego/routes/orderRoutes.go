package routes

import (
	"bego/controllers"
	"bego/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterOrderRoutes(router *mux.Router) {
	// POST /api/orders - user-only (mock payment + order placement)
	router.Handle("/api/orders", middleware.ProtectRoutes(http.HandlerFunc(controllers.CreateOrder))).Methods("POST")

	// GET /api/orders - user-only (view order history)
	router.Handle("/api/orders", middleware.ProtectRoutes(http.HandlerFunc(controllers.GetOrderHistory))).Methods("GET")
}
