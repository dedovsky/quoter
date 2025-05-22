package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"quoter/models"
	"quoter/storage"
	"strconv"
)

type QuoteHandler struct {
	store storage.Store
}

func NewQuoteHandler(store storage.Store) *QuoteHandler {
	return &QuoteHandler{store: store}
}

func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var quote models.Quote

	if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
		http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	if quote.Author == "" || quote.Quote == "" {
		http.Error(w, "Author или quote не может быть пустым", http.StatusBadRequest)
		return
	}
	id := h.store.AddQuote(quote)

	quote.ID = id
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(quote)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *QuoteHandler) GetAllQuotes(w http.ResponseWriter, r *http.Request) {
	quotes := h.store.GetAllQuotes()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(quotes)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quotes := h.store.GetAllQuotes()
	quotesCount := len(quotes)

	if quotesCount == 0 {
		http.Error(w, "Цитат нет", http.StatusNotFound)
		return
	}

	randomIndex := rand.Intn(quotesCount)
	randomQuote := quotes[randomIndex]
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(randomQuote)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *QuoteHandler) GetQuotesByAuthor(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	if author == "" {
		http.Error(w, "Параметр author не может быть пустым", http.StatusBadRequest)
		return
	}

	quotes := h.store.GetQuotesByAuthor(author)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(quotes)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *QuoteHandler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if !h.store.DeleteQuote(id) {
		http.Error(w, "Цитата не найдена", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
