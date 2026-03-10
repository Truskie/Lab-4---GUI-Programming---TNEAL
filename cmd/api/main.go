package main

import (
	"log"
	"net/http"
	"os"

	"lab2-terrylneal/internal/db"
	"lab2-terrylneal/internal/routes"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgresql://lab2_user:lab2pass@localhost/lab2_terrylneal?sslmode=disable"
		log.Println("DB_URL not set, using default connection string")
	}

	database, err := db.NewDB(dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer database.Close()

	log.Println("Successfully connected to database")

	mux := http.NewServeMux()

	routes.SetupRoutes(mux, database)

	handler := routes.ApplyMiddleware(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.Printf("Server Starting on :%s", port)
	err = http.ListenAndServe(":"+port, handler)
	log.Fatal(err)
}
