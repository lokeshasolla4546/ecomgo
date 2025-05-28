package main

import (
	"bego/config"
	"bego/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.Connect()
	router := mux.NewRouter()
	routes.RegisterUserRoutes(router)
	routes.RegisterProductRoutes(router)
	routes.RegisterCartRoutes(router)
	routes.RegisterOrderRoutes(router)
	log.Println("Server running at http://localhost:5000")
	err := http.ListenAndServe(":5000", router)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
