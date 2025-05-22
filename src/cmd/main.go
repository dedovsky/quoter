package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"quoter/handlers"
	"quoter/storage"
)

func main() {
	r := mux.NewRouter()
	store := storage.NewMemoryStorage()
	handler := handlers.NewQuoteHandler(store)

	r.HandleFunc("/quotes", handler.CreateQuote).Methods("POST")
	r.HandleFunc("/quotes", handler.GetAllQuotes).Methods("GET")
	r.HandleFunc("/quotes/random", handler.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes", handler.GetQuotesByAuthor).Queries("author", "{author}").Methods("GET")
	r.HandleFunc("/quotes/{id}", handler.DeleteQuote).Methods("DELETE")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
