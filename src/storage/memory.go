package storage

import (
	"quoter/models"
	"sync"
)

type Store interface {
	AddQuote(quote models.Quote) int
	GetAllQuotes() []models.Quote
	GetQuotesByAuthor(author string) []models.Quote
	DeleteQuote(id int) bool
}
type MemoryStore struct {
	quotes []models.Quote
	mu     sync.Mutex
	nextID int
}

func NewMemoryStorage() *MemoryStore {
	return &MemoryStore{
		nextID: 1,
		quotes: []models.Quote{},
	}
}

func (s *MemoryStore) AddQuote(quote models.Quote) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	quote.ID = s.nextID
	s.nextID++

	s.quotes = append(s.quotes, quote)

	return quote.ID
}

func (s *MemoryStore) GetAllQuotes() []models.Quote {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.quotes
}

func (s *MemoryStore) GetQuotesByAuthor(author string) []models.Quote {
	s.mu.Lock()
	defer s.mu.Unlock()

	var authorQuotes []models.Quote

	for _, quote := range s.quotes {
		if quote.Author == author {
			authorQuotes = append(authorQuotes, quote)
		}
	}

	return authorQuotes
}

func (s *MemoryStore) DeleteQuote(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, quote := range s.quotes {
		if quote.ID == id {
			s.quotes = append(s.quotes[:i], s.quotes[i+1:]...)
			return true
		}
	}

	return false
}
