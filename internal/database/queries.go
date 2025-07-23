package database

import (
	"context"
	"database/sql"
	"digicert-library-app/internal/models"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func (d *Database) GetBooks(ctx context.Context, limit, offset int) ([]models.Book, error) {
	// move DB queries to DB folders
	books := []models.Book{}
	query := "SELECT id, title FROM books LIMIT ? OFFSET ?"
	rows, err := d.Conn.QueryContext(context.Background(), query, limit, offset)
	if err != nil {
		return books, err
	}
	defer rows.Close()
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title); err != nil {
			return books, err
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return books, err
	}

	return books, err
}

func (d *Database) GetBookByID(ctx context.Context, id string) (models.Book, error) {
	query := "select * from books where id = ?"
	row := d.Conn.QueryRowContext(context.Background(), query, id)
	var book models.Book
	if err := row.Scan(&book.ID, &book.Title); err != nil {
		if err == sql.ErrNoRows {
			return book, err
		}
		return book, err
	}
	return book, nil
}

func (d *Database) CreateBook(ctx context.Context, newBook models.Book) (string, error) {
	id := uuid.New() // Generate a new UUID for the book
	query := "INSERT INTO books (id, title) VALUES (?, ?)"
	_, err := d.Conn.ExecContext(context.Background(), query, id, newBook.Title)
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
	query := "UPDATE books SET title = ? WHERE id = ?"
	result, err := d.Conn.ExecContext(context.Background(), query, updatedBook.Title, id)
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
	result, err := d.Conn.ExecContext(context.Background(), query, id)
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
