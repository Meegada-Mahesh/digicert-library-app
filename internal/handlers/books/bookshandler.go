package books

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"digicert-library-app/internal/database"
	"digicert-library-app/internal/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type BooksHandler struct {
	db *database.Database
}

func InitBooksHandler(ctx context.Context, db *database.Database) *BooksHandler {
	// Initialize the Book handler
	return &BooksHandler{
		db: db,
	}
}
func (b *BooksHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")
	ctx := r.Context()
	page := 1
	limit := 10
	if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
		limit = l
	}
	offset := (page - 1) * limit

	books, err := b.db.GetBooks(ctx, limit, offset)
	if err != nil {
		http.Error(w, "Error in fetching books from library", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"data":  books,
		"page":  page,
		"limit": limit,
	}
	json.NewEncoder(w).Encode(response)
}

func (b *BooksHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	ctx := r.Context()
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid book ID format", http.StatusBadRequest)
		return
	}

	book, err := b.db.GetBookByID(ctx, id)
	if err != nil {
		http.Error(w, "Error in fetching books from library", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"book": book,
	}
	json.NewEncoder(w).Encode(response)
}

func (b *BooksHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if newBook.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	_, err := b.db.CreateBook(r.Context(), newBook)
	if err != nil {
		http.Error(w, "Couldn't create book", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Book created"})
}

func (b *BooksHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updateBook); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid book ID format", http.StatusBadRequest)
		return
	}

	resultMsg, err := b.db.UpdateBook(r.Context(), id, updateBook)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": resultMsg})
}

func (b *BooksHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid book ID format", http.StatusBadRequest)
		return
	}

	resultMsg, err := b.db.DeleteBook(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": resultMsg})
}
