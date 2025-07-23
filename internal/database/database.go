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

	// Move DB logic to internal folder

	_ = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	fmt.Println("Database connection established successfully")

	return db, nil
}
