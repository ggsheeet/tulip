package api

import (
	"github.com/ggsheet/kerigma/internal/database"
)

type APIError struct {
	Error string `json:"message"`
}

type APIServer struct {
	account  *AccountHandlers
	book     *BookHandlers
	article  *ArticleHandlers
	resource *ResourceHandlers
	order    *OrderHandlers
}

type AccountHandlers struct {
	db database.AccountInterface
}

type BookHandlers struct {
	db database.BookInterface
}

type ArticleHandlers struct {
	db database.ArticleInterface
}

type ResourceHandlers struct {
	db database.ResourceInterface
}

type OrderHandlers struct {
	db database.OrderInterface
}
