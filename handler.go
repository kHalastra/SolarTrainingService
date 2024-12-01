package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book struct to model the data
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// BookHandler struct that will hold the DB connection and handlers
type BookHandler struct {
	DB *sql.DB
}

// CreateBook godoc
// @Summary Create a new book
// @Description Add a new book to the database
// @Tags books
// @Accept json
// @Produce json
// @Param book body Book true "Book data"
// @Success 200 {object} Book
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Server error"
// @Router /books [post]
func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	err := h.DB.QueryRow("INSERT INTO books (title, author, year) VALUES ($1, $2, $3) RETURNING id",
		book.Title, book.Author, book.Year).Scan(&book.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// GetBooks godoc
// @Summary List all books
// @Description Get all books from the database
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} Book
// @Failure 500 {string} string "Server error"
// @Router /books [get]
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, title, author, year FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

// GetBook godoc
// @Summary Get a book by ID
// @Description Retrieve a specific book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} Book
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Book not found"
// @Failure 500 {string} string "Server error"
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book Book
	err = h.DB.QueryRow("SELECT id, title, author, year FROM books WHERE id = $1", id).
		Scan(&book.ID, &book.Title, &book.Author, &book.Year)

	if err == sql.ErrNoRows {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// UpdateBook godoc
// @Summary Update a book by ID
// @Description Update the details of an existing book
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body Book true "Updated book data"
// @Success 200 {string} string "Book updated successfully"
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "Book not found"
// @Failure 500 {string} string "Server error"
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	_, err = h.DB.Exec("UPDATE books SET title = $1, author = $2, year = $3 WHERE id = $4",
		book.Title, book.Author, book.Year, id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book updated successfully")
}

// DeleteBook godoc
// @Summary Delete a book by ID
// @Description Remove a book from the database by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {string} string "Book deleted successfully"
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Book not found"
// @Failure 500 {string} string "Server error"
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Book deleted successfully")
}
