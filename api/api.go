package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ggsheet/kerigma/internal/database"
)

func NewAPIServer(listenAddr string, accInterface database.AccountInterface, bookInterface database.BookInterface) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		account:    &AccountHandlers{accountInterface: accInterface},
		book:       &BookHandlers{bookInterface: bookInterface},
	}
}

func (s *APIServer) Run() {
	r := http.NewServeMux()

	r.HandleFunc("/account", makeHTTPHandleFunc(s.account.handleAccount))
	r.HandleFunc("GET /account/{id}", makeHTTPHandleFunc(s.account.handleGetAccountById))
	r.HandleFunc("DELETE /account/{id}", makeHTTPHandleFunc(s.account.handleDeleteAccount))

	r.HandleFunc("/book", makeHTTPHandleFunc(s.book.handleBook))
	r.HandleFunc("GET /book/{id}", makeHTTPHandleFunc(s.book.handleGetBookById))
	r.HandleFunc("DELETE /book/{id}", makeHTTPHandleFunc(s.book.handleDeleteBook))

	log.Println("JSON API server on port: :8080")

	http.ListenAndServe(":8080", r)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
