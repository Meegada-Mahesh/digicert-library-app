package database

import (
	"database/sql"
	"digicert-library-app/internal/models"
	"strings"

	"github.com/google/uuid"
)

// Dynamic placeholder generation function
func generatePlaceholders(count int) string {
	if count <= 0 {
		return ""
	}
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}

// Column management functions
func getBookColumns() []string {
	return []string{"id", "title", "author", "published_year", "genre"}
}

func getInsertColumns() []string {
	return []string{"id", "title", "author", "published_year", "genre"}
}

func getUpdateColumns() []string {
	return []string{"title", "author", "published_year", "genre"}
}

// Helper functions using dynamic placeholders
func getBookColumnsString() string {
	return strings.Join(getBookColumns(), ", ")
}

func getInsertColumnsString() string {
	return strings.Join(getInsertColumns(), ", ")
}

func getInsertPlaceholders() string {
	return generatePlaceholders(len(getInsertColumns()))
}

func getUpdateColumnsString() string {
	columns := getUpdateColumns()
	updateParts := make([]string, len(columns))
	for i, col := range columns {
		updateParts[i] = col + " = ?"
	}
	return strings.Join(updateParts, ", ")
}

func scanBookRow(row *sql.Row) (models.Book, error) {
	var book models.Book
	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear, &book.Genre)
	return book, err
}

func scanBookRows(rows *sql.Rows) (models.Book, error) {
	var book models.Book
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishedYear, &book.Genre)
	return book, err
}

func getInsertValues(id uuid.UUID, book models.Book) []interface{} {
	return []interface{}{id, book.Title, book.Author, book.PublishedYear, book.Genre}
}

func getUpdateValues(book models.Book, id string) []interface{} {
	return []interface{}{book.Title, book.Author, book.PublishedYear, book.Genre, id}
}
