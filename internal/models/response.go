package models

type BookResponse struct {
	Book Book `json:"book"`
}

type BooksResponse struct {
	Data       []Book `json:"data"`
	Page       int    `json:"page,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Total      int    `json:"total,omitempty"`
	TotalPages int    `json:"totalPages,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
