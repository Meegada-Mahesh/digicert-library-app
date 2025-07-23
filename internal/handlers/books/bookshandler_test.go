package books

import (
	"bytes"
	"context"
	"digicert-library-app/internal/database"
	"digicert-library-app/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

func getMockHandler(t *testing.T) (*BooksHandler, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	handler := InitBooksHandler(context.Background(), &database.Database{Conn: db})
	return handler, mock
}

func TestGetBooks_WithData(t *testing.T) {
	handler, mock := getMockHandler(t)
	rows := sqlmock.NewRows([]string{"id", "title", "author", "published_year", "genre"}).
		AddRow("123e4567-e89b-12d3-a456-426614174000", "Book One", "Author One", 2020, "Fiction").
		AddRow("123e4567-e89b-12d3-a456-426614174001", "Book Two", "Author Two", 2021, "Non-Fiction")
	mock.ExpectQuery("SELECT id, title, author, published_year, genre FROM books").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	handler.GetBooks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var resp models.BooksResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Errorf("Expected 2 books, got %d", len(resp.Data))
	}
}

func TestGetBooks_Empty(t *testing.T) {
	handler, mock := getMockHandler(t)
	rows := sqlmock.NewRows([]string{"id", "title", "author", "published_year", "genre"})
	mock.ExpectQuery("SELECT id, title, author, published_year, genre FROM books").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	handler.GetBooks(w, req)

	if w.Code != http.StatusNoContent && w.Code != http.StatusOK {
		t.Errorf("Expected status 204 or 200, got %d", w.Code)
	}
}

func TestGetBooks_DBError(t *testing.T) {
	handler, mock := getMockHandler(t)
	mock.ExpectQuery("SELECT id, title, author, published_year, genre FROM books").WillReturnError(errors.New("db error"))

	req := httptest.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	handler.GetBooks(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
	var resp models.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

func TestGetBookByID_Valid(t *testing.T) {
	handler, mock := getMockHandler(t)
	rows := sqlmock.NewRows([]string{"id", "title", "author", "published_year", "genre"}).
		AddRow("123e4567-e89b-12d3-a456-426614174000", "Test Book", "Test Author", 2023, "Fiction")
	mock.ExpectQuery("SELECT id, title, author, published_year, genre FROM books WHERE id = \\?").
		WithArgs("123e4567-e89b-12d3-a456-426614174000").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/books/123e4567-e89b-12d3-a456-426614174000", nil)
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})
	handler.GetBookByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var resp models.BookResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Book.Title == "" {
		t.Errorf("Expected book title, got empty string")
	}
}

func TestCreateBook_ValidPayload(t *testing.T) {
	handler, mock := getMockHandler(t)
	book := models.Book{Title: "Test Book", Author: "Test Author", PublishedYear: 2023, Genre: "Fiction"}
	body, _ := json.Marshal(book)
	mock.ExpectExec("INSERT INTO books").
		WithArgs(sqlmock.AnyArg(), "Test Book", "Test Author", 2023, "Fiction").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.CreateBook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var resp models.MessageResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Message == "" {
		t.Errorf("Expected message, got empty string")
	}
}

func TestCreateBook_MissingTitle(t *testing.T) {
	handler, _ := getMockHandler(t)
	book := models.Book{Author: "Test Author"}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.CreateBook(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
	var resp models.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

func TestCreateBook_MissingAuthor(t *testing.T) {
	handler, _ := getMockHandler(t)
	book := models.Book{Title: "Test Book"}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.CreateBook(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
	var resp models.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

func TestCreateBook_InvalidPayload(t *testing.T) {
	handler, _ := getMockHandler(t)
	req := httptest.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`invalid-json`)))
	w := httptest.NewRecorder()
	handler.CreateBook(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
	var resp models.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

func TestCreateBook_DBError(t *testing.T) {
	handler, mock := getMockHandler(t)
	book := models.Book{Title: "Test Book", Author: "Test Author"}
	body, _ := json.Marshal(book)
	mock.ExpectExec("INSERT INTO books").
		WithArgs(sqlmock.AnyArg(), "Test Book", "Test Author", 0, "").
		WillReturnError(errors.New("insert error"))

	req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.CreateBook(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
	var resp models.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

func TestUpdateBook_Valid(t *testing.T) {
	handler, mock := getMockHandler(t)
	book := models.Book{Title: "Updated Title", Author: "Updated Author", PublishedYear: 2024, Genre: "Updated Genre"}
	body, _ := json.Marshal(book)
	mock.ExpectExec("UPDATE books SET title = \\?, author = \\?, published_year = \\?, genre = \\? WHERE id = \\?").
		WithArgs("Updated Title", "Updated Author", 2024, "Updated Genre", "123e4567-e89b-12d3-a456-426614174000").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest("PUT", "/books/123e4567-e89b-12d3-a456-426614174000", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})
	handler.UpdateBook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var resp models.MessageResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Message == "" {
		t.Errorf("Expected message, got empty string")
	}
}

func TestUpdateBook_NotFound(t *testing.T) {
	handler, mock := getMockHandler(t)
	book := models.Book{Title: "Updated Title", Author: "Updated Author"}
	body, _ := json.Marshal(book)
	mock.ExpectExec("UPDATE books SET title = \\?, author = \\?, published_year = \\?, genre = \\? WHERE id = \\?").
		WithArgs("Updated Title", "Updated Author", 0, "", "123e4567-e89b-12d3-a456-426614174000").
		WillReturnResult(sqlmock.NewResult(1, 0))

	req := httptest.NewRequest("PUT", "/books/123e4567-e89b-12d3-a456-426614174000", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})
	handler.UpdateBook(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
	var resp models.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

func TestUpdateBook_InvalidPayload(t *testing.T) {
	handler, _ := getMockHandler(t)
	req := httptest.NewRequest("PUT", "/books/123e4567-e89b-12d3-a456-426614174000", bytes.NewBuffer([]byte(`invalid-json`)))
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})
	handler.UpdateBook(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
	var resp models.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

func TestDeleteBook_Valid(t *testing.T) {
	handler, mock := getMockHandler(t)
	mock.ExpectExec("DELETE FROM books WHERE id = \\?").
		WithArgs("123e4567-e89b-12d3-a456-426614174000").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest("DELETE", "/books/123e4567-e89b-12d3-a456-426614174000", nil)
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})
	handler.DeleteBook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var resp models.MessageResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Message == "" {
		t.Errorf("Expected message, got empty string")
	}
}

func TestDeleteBook_NotFound(t *testing.T) {
	handler, mock := getMockHandler(t)
	mock.ExpectExec("DELETE FROM books WHERE id = \\?").
		WithArgs("123e4567-e89b-12d3-a456-426614174000").
		WillReturnResult(sqlmock.NewResult(1, 0))

	req := httptest.NewRequest("DELETE", "/books/123e4567-e89b-12d3-a456-426614174000", nil)
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})
	handler.DeleteBook(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
	var resp models.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

func TestDeleteBook_DBError(t *testing.T) {
	handler, mock := getMockHandler(t)
	mock.ExpectExec("DELETE FROM books WHERE id = \\?").
		WithArgs("123e4567-e89b-12d3-a456-426614174000").
		WillReturnError(errors.New("delete error"))

	req := httptest.NewRequest("DELETE", "/books/123e4567-e89b-12d3-a456-426614174000", nil)
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})
	handler.DeleteBook(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
	var resp models.ErrorResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Error == "" {
		t.Errorf("Expected error message, got empty string")
	}
}

// Helper to set mux vars for path parameters in tests
func muxSetVars(r *http.Request, vars map[string]string) *http.Request {
	return mux.SetURLVars(r, vars)
}
