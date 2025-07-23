package main

import (
	"context"
	"database/sql"
	"digicert-library-app/internal/database"
	"digicert-library-app/internal/middleware"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pressly/goose/v3"

	"digicert-library-app/internal/handlers/books"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

func waitForDatabase(db *sql.DB, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err == nil {
			log.Println("Database connection established")
			return nil
		}
		log.Printf("Waiting for database... attempt %d/%d", i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
}

func main() {

	// create a database connection
	ctx := context.Background()
	db, err := database.NewDBConnection()
	if err != nil {
		log.Fatalf("Error in DB connection error: %v", err)
	}
	defer db.Conn.Close()

	// Set connection pool settings to prevent connection drops
	db.Conn.SetMaxOpenConns(25)
	db.Conn.SetMaxIdleConns(25)
	db.Conn.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connection created, testing connectivity...")

	// Wait for database to be ready with better error handling
	if err := waitForDatabase(db.Conn, 30); err != nil {
		log.Printf("Database connection test failed: %v", err)
		log.Println("Attempting to reconnect...")

		// Try to reconnect once more
		db, err = database.NewDBConnection()
		if err != nil {
			log.Fatalf("Database reconnection failed: %v", err)
		}

		if err := waitForDatabase(db.Conn, 10); err != nil {
			log.Fatalf("Database not ready after reconnection: %v", err)
		}
	}

	//setting up goose
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("mysql"); err != nil {
		log.Fatal("error in setting mysql dialect", err)
	}

	if err := goose.Up(db.Conn, "db/migrations"); err != nil {
		log.Fatal("error in setting up migrations", err)
	}

	// initialize the books handler
	booksHandler := books.InitBooksHandler(ctx, db)

	// routing logic
	r := mux.NewRouter()
	// Set up CORS middleware
	r.Use(mux.CORSMethodMiddleware(r))

	// Adding middlewares
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.JsonHeaderMiddleware)
	r.Use(middleware.LimitBodySizeMiddleware)
	r.Use(middleware.AuthMiddleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to Digicert!")
	})

	r.HandleFunc("/books", booksHandler.GetBooks).Methods("GET")
	r.HandleFunc("/books/{id}", booksHandler.GetBookByID).Methods("GET")
	r.HandleFunc("/books", booksHandler.CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", booksHandler.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", booksHandler.DeleteBook).Methods("DELETE")

	srv := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
