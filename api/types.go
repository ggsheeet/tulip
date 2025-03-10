package api

import (
	"github.com/ggsheet/tulip/internal/database"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/resend/resend-go/v2"
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
	payment  *MPServer
	mailing  *ResendServer
}

type MPServer struct {
	cfg *config.Config
	a   database.AccountInterface
	b   database.BookInterface
	o   database.OrderInterface
	m   *ResendServer
}

type ResendServer struct {
	msg *resend.Client
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
