package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	var err error
	connStr := "host=localhost user=asolla password=1234 dbname=ironman sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("DB not reachable:", err)
	}
	fmt.Println(" Connected to PostgreSQL database: ironman")
}
