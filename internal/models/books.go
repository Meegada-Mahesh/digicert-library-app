package models

import "github.com/google/uuid"

type Book struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	PublishedYear int       `json:"published_year,omitempty"`
	Genre         string    `json:"genre,omitempty"`
}
