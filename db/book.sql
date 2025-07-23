-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS books (
    id CHAR(36) NOT NULL PRIMARY KEY,           -- UUID for each book
    title VARCHAR(255) NOT NULL,                -- Book title
    author VARCHAR(255) NOT NULL,               -- Author name
    published_year INT,                         -- Year of publication
    genre VARCHAR(100),                         -- Book genre/category
    INDEX idx_title (title),                    -- Index for faster title searches
    INDEX idx_author (author)                   -- Index for faster author searches
);

-- +goose StatementEnd
-- +goose Down