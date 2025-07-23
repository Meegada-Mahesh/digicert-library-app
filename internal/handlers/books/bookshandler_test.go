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
	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow("123e4567-e89b-12d3-a456-426614174000", "Book One").
		AddRow("123e4567-e89b-12d3-a456-426614174001", "Book Two")
	mock.ExpectQuery("SELECT id, title FROM books").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	handler.GetBooks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var resp map[string][]models.Book
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	if len(resp["data"]) != 2 {
		t.Errorf("Expected 2 books, got %d", len(resp["data"]))
	}
}

func TestGetBooks_Empty(t *testing.T) {
	handler, mock := getMockHandler(t)
	rows := sqlmock.NewRows([]string{"id", "title"})
	mock.ExpectQuery("SELECT id, title FROM books").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	handler.GetBooks(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}
}

func TestGetBooks_DBError(t *testing.T) {
	handler, mock := getMockHandler(t)
	mock.ExpectQuery("SELECT id, title FROM books").WillReturnError(errors.New("db error"))

	req := httptest.NewRequest("GET", "/books", nil)
	w := httptest.NewRecorder()
	handler.GetBooks(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestCreateBook_ValidPayload(t *testing.T) {
	handler, mock := getMockHandler(t)
	book := models.Book{Title: "Test Book"}
	body, _ := json.Marshal(book)
	mock.ExpectExec("INSERT INTO books").
		WithArgs(sqlmock.AnyArg(), "Test Book").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.CreateBook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
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
}

func TestCreateBook_DBError(t *testing.T) {
	handler, mock := getMockHandler(t)
	book := models.Book{Title: "Test Book"}
	body, _ := json.Marshal(book)
	mock.ExpectExec("INSERT INTO books").
		WithArgs(sqlmock.AnyArg(), "Test Book").
		WillReturnError(errors.New("insert error"))

	req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	handler.CreateBook(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestUpdateBook_Valid(t *testing.T) {
	handler, mock := getMockHandler(t)
	book := models.Book{Title: "Updated Title"}
	body, _ := json.Marshal(book)
	mock.ExpectExec("UPDATE books SET title = \\? WHERE id = \\?").
		WithArgs("Updated Title", "123e4567-e89b-12d3-a456-426614174000").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest("PUT", "/books/123e4567-e89b-12d3-a456-426614174000", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "123e4567-e89b-12d3-a456-426614174000"})
	handler.UpdateBook(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestUpdateBook_NotFound(t *testing.T) {
	handler, mock := getMockHandler(t)
	book := models.Book{Title: "Updated Title"}
	body, _ := json.Marshal(book)
	mock.ExpectExec("UPDATE books SET title = \\? WHERE id = \\?").
		WithArgs("Updated Title", "unknown-id").
		WillReturnResult(sqlmock.NewResult(1, 0))

	req := httptest.NewRequest("PUT", "/books/unknown-id", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "unknown-id"})
	handler.UpdateBook(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUpdateBook_InvalidPayload(t *testing.T) {
	handler, _ := getMockHandler(t)
	req := httptest.NewRequest("PUT", "/books/123", bytes.NewBuffer([]byte(`invalid-json`)))
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "123"})
	handler.UpdateBook(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
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
}

func TestDeleteBook_NotFound(t *testing.T) {
	handler, mock := getMockHandler(t)
	mock.ExpectExec("DELETE FROM books WHERE id = \\?").
		WithArgs("unknown-id").
		WillReturnResult(sqlmock.NewResult(1, 0))

	req := httptest.NewRequest("DELETE", "/books/unknown-id", nil)
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "unknown-id"})
	handler.DeleteBook(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeleteBook_DBError(t *testing.T) {
	handler, mock := getMockHandler(t)
	mock.ExpectExec("DELETE FROM books WHERE id = \\?").
		WithArgs("123").
		WillReturnError(errors.New("delete error"))

	req := httptest.NewRequest("DELETE", "/books/123", nil)
	w := httptest.NewRecorder()
	req = muxSetVars(req, map[string]string{"id": "123"})
	handler.DeleteBook(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

// Helper to set mux vars for path parameters in tests
func muxSetVars(r *http.Request, vars map[string]string) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), "vars", vars))
}
