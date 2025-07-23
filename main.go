package main

import (
	"context"
	"digicert-library-app/internal/database"
	"digicert-library-app/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"digicert-library-app/internal/handlers/books"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	// create a database connection
	ctx := context.Background()
	db, err := database.NewDBConnection()
	if err != nil {
		log.Fatalf("Error in DB connection error: %v", err)
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

// Write a READMe file
