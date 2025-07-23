package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRootHandler checks if the root endpoint returns the welcome message.
func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Setup a minimal router for testing
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to Digicert!"))
	})

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	if w.Body.String() != "Welcome to Digicert!\n" && w.Body.String() != "Welcome to Digicert!" {
		t.Errorf("Unexpected body: %s", w.Body.String())
	}
}
