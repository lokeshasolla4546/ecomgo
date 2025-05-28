package routes

import (
	"bego/controllers"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/api/users/register", controllers.Register).Methods("POST")
	router.HandleFunc("/api/users/login", controllers.Login).Methods("POST")
}
