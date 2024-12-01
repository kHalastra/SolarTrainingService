package main

import (
	_ "SolarTrainingService/docs"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

var db *sql.DB

func initDB(config *Config) (*sql.DB, error) {
	var err error

	psqlUrl := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		5432,
		config.DBName)

	db, err = sql.Open("postgres", psqlUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to ensure a proper connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database successfully!")
	return db, nil
}

func main() {
	config := LoadConfig()

	// Initialize DB connection
	db, err := initDB(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize the BookHandler with the database connection
	handler := &BookHandler{DB: db}

	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// Define routes using the handler's methods
	r.HandleFunc("/books", handler.CreateBook).Methods("POST")
	r.HandleFunc("/books", handler.GetBooks).Methods("GET")
	r.HandleFunc("/books/{id}", handler.GetBook).Methods("GET")
	r.HandleFunc("/books/{id}", handler.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", handler.DeleteBook).Methods("DELETE")

	// Add Swagger endpoint
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start the server
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", r)
}
