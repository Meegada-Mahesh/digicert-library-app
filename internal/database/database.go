package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Database struct {
	Conn *sql.DB
}

func NewDBConnection() (*Database, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}

	return &Database{Conn: db}, nil
}

func connect() (*sql.DB, error) {
	_ = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbName)
	log.Printf("Connecting to database with DSN: %s:****@tcp(%s:3306)/%s", dbUser, dbHost, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Database connection failed: %v", err)
		return nil, err
	}

	// Test the connection immediately
	err = db.Ping()
	if err != nil {
		log.Printf("Database ping failed: %v", err)
		return nil, err
	}

	log.Println("Database connection and ping successful")
	return db, nil
}
