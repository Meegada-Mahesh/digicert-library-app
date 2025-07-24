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
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")

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
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Failed to fetch books"})
		return
	}

	json.NewEncoder(w).Encode(models.BooksResponse{
		Data:  books,
		Page:  page,
		Limit: limit,
	})
}

func (b *BooksHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars["id"]
	if _, err := uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid book ID format"})
		return
	}

	book, err := b.db.GetBookByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Book not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Error in fetching books from library"})
		return
	}
	json.NewEncoder(w).Encode(models.BookResponse{Book: book})
}

func (b *BooksHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	var newBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid request payload"})
		return
	}

	// Validate required fields
	if newBook.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Title is required"})
		return
	}
	if newBook.Author == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Author is required"})
		return
	}

	_, err := b.db.CreateBook(ctx, newBook)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Couldn't create book"})
		return
	}
	json.NewEncoder(w).Encode(models.MessageResponse{Message: "Book created"})
}

func (b *BooksHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	var updateBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&updateBook); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid request payload"})
		return
	}
	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid book ID format"})
		return
	}

	// Validate required fields for update
	if updateBook.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Title is required"})
		return
	}
	if updateBook.Author == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Author is required"})
		return
	}

	resultMsg, err := b.db.UpdateBook(ctx, id, updateBook)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Book not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Failed to update book"})
		return
	}
	json.NewEncoder(w).Encode(models.MessageResponse{Message: resultMsg})
}

func (b *BooksHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars["id"]

	if _, err := uuid.Parse(id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Invalid book ID format"})
		return
	}

	resultMsg, err := b.db.DeleteBook(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Book not found"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "Failed to delete book"})
		return
	}
	json.NewEncoder(w).Encode(models.MessageResponse{Message: resultMsg})
}
