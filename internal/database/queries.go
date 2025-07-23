package database

import (
	"context"
	"database/sql"
	"digicert-library-app/internal/models"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// CRUD functions using the helper functions
func (d *Database) GetBooks(ctx context.Context, limit, offset int) ([]models.Book, error) {
	books := []models.Book{}
	query := "SELECT " + getBookColumnsString() + " FROM books LIMIT ? OFFSET ?"
	rows, err := d.Conn.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return books, err
	}
	defer rows.Close()

	for rows.Next() {
		book, err := scanBookRows(rows)
		if err != nil {
			return books, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return books, err
	}
	return books, nil
}

func (d *Database) GetBookByID(ctx context.Context, id string) (models.Book, error) {
	query := "SELECT " + getBookColumnsString() + " FROM books WHERE id = ?"
	row := d.Conn.QueryRowContext(ctx, query, id)

	book, err := scanBookRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return book, err
		}
		return book, err
	}
	return book, nil
}

func (d *Database) CreateBook(ctx context.Context, newBook models.Book) (string, error) {
	id := uuid.New()
	query := "INSERT INTO books (" + getInsertColumnsString() + ") VALUES (" + getInsertPlaceholders() + ")"
	values := getInsertValues(id, newBook)

	_, err := d.Conn.ExecContext(ctx, query, values...)
	if err != nil {
		// Check for duplicate entry error (MySQL error code 1062)
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return "", err
		}
		return "", err
	}
	return "Book Inserted", nil
}

func (d *Database) UpdateBook(ctx context.Context, id string, updatedBook models.Book) (string, error) {
	query := "UPDATE books SET " + getUpdateColumnsString() + " WHERE id = ?"
	values := getUpdateValues(updatedBook, id)

	result, err := d.Conn.ExecContext(ctx, query, values...)
	if err != nil {
		return "", err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if rowsAffected == 0 {
		return "", sql.ErrNoRows
	}
	return "Book Updated", nil
}

func (d *Database) DeleteBook(ctx context.Context, id string) (string, error) {
	query := "DELETE FROM books WHERE id = ?"
	result, err := d.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return "", err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	if rowsAffected == 0 {
		return "", sql.ErrNoRows
	}
	return "Book Deleted", nil
}
