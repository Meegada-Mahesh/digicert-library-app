package models

import "github.com/google/uuid"

type Book struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}
