package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"quoter/models"
	"quoter/storage"
	"testing"

	"github.com/gorilla/mux"
)

func setupTestRouter() (*mux.Router, *QuoteHandler) {
	store := storage.NewMemoryStorage()
	handler := NewQuoteHandler(store)
	router := mux.NewRouter()
	return router, handler
}

func setupTestRouterWithQuotes(quotes []models.Quote) (*mux.Router, *QuoteHandler) {
	store := storage.NewMemoryStorage()
	for _, q := range quotes {
		store.AddQuote(q)
	}
	handler := NewQuoteHandler(store)
	router := mux.NewRouter()
	return router, handler
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("Expected status %d, got %d", want, got)
	}
}

func TestCreateQuote(t *testing.T) {
	router, handler := setupTestRouter()
	router.HandleFunc("/quotes", handler.CreateQuote).Methods("POST")

	quote := models.Quote{Author: "Confucius", Quote: "Life is simple, but we insist on making it complicated."}
	body, _ := json.Marshal(quote)
	req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assertStatus(t, w.Code, http.StatusCreated)

	var response models.Quote
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if response.ID != 1 {
		t.Errorf("Expected ID 1, got %d", response.ID)
	}
}

func TestGetAllQuotes(t *testing.T) {
	quotes := []models.Quote{
		{ID: 1, Author: "Confucius", Quote: "Life is simple."},
	}
	router, handler := setupTestRouterWithQuotes(quotes)
	router.HandleFunc("/quotes", handler.GetAllQuotes).Methods("GET")

	req := httptest.NewRequest("GET", "/quotes", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assertStatus(t, w.Code, http.StatusOK)

	var gotQuotes []models.Quote
	if err := json.NewDecoder(w.Body).Decode(&gotQuotes); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(gotQuotes) != len(quotes) {
		t.Errorf("Expected %d quote(s), got %d", len(quotes), len(gotQuotes))
	}
}

func TestGetRandomQuote(t *testing.T) {
	quotes := []models.Quote{
		{ID: 1, Author: "Confucius", Quote: "Life is simple."},
	}
	router, handler := setupTestRouterWithQuotes(quotes)
	router.HandleFunc("/quotes/random", handler.GetRandomQuote).Methods("GET")

	req := httptest.NewRequest("GET", "/quotes/random", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assertStatus(t, w.Code, http.StatusOK)

	var quote models.Quote
	if err := json.NewDecoder(w.Body).Decode(&quote); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if quote.ID != 1 {
		t.Errorf("Expected ID 1, got %d", quote.ID)
	}
}

func TestGetQuotesByAuthor(t *testing.T) {
	quotes := []models.Quote{
		{ID: 1, Author: "Confucius", Quote: "Life is simple."},
		{ID: 2, Author: "Plato", Quote: "The greatest wealth is to live content with little."},
	}
	router, handler := setupTestRouterWithQuotes(quotes)
	router.HandleFunc("/quotes", handler.GetQuotesByAuthor).Queries("author", "{author}").Methods("GET")

	req := httptest.NewRequest("GET", "/quotes?author=Confucius", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assertStatus(t, w.Code, http.StatusOK)

	var filtered []models.Quote
	if err := json.NewDecoder(w.Body).Decode(&filtered); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(filtered) != 1 {
		t.Errorf("Expected 1 quote, got %d", len(filtered))
	}
}

func TestDeleteQuote(t *testing.T) {
	quotes := []models.Quote{
		{ID: 1, Author: "Confucius", Quote: "Life is simple."},
	}
	router, handler := setupTestRouterWithQuotes(quotes)
	router.HandleFunc("/quotes/{id}", handler.DeleteQuote).Methods("DELETE")
	router.HandleFunc("/quotes", handler.GetAllQuotes).Methods("GET")

	req := httptest.NewRequest("DELETE", "/quotes/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assertStatus(t, w.Code, http.StatusNoContent)

	req = httptest.NewRequest("GET", "/quotes", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var remaining []models.Quote
	if err := json.NewDecoder(w.Body).Decode(&remaining); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(remaining) != 0 {
		t.Errorf("Expected 0 quotes, got %d", len(remaining))
	}
}
