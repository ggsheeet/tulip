package api

import (
	"net/http"

	"github.com/ggsheet/kerigma/internal/database"
)

type APIServer struct {
	listenAddr string
	account    *AccountHandlers
	book       *BookHandlers
}

type APIFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string
}

type AccountHandlers struct {
	accountInterface database.AccountInterface
}

type BookHandlers struct {
	bookInterface database.BookInterface
}
