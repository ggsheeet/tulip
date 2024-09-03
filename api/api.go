package api

import (
	"github.com/ggsheet/kerigma/internal/database"
	"github.com/labstack/echo/v4"
)

func NewAPIServer(db *database.PostgresDB) *APIServer {
	return &APIServer{
		account:  &AccountHandlers{db},
		book:     &BookHandlers{db},
		article:  &ArticleHandlers{db},
		resource: &ResourceHandlers{db},
		order:    &OrderHandlers{db},
	}
}

func (s *APIServer) APIRouter(e *echo.Echo) {
	apiGroup := e.Group("/api")
	apiGroup.Use(authMiddleware)

	apiGroup.Any("/account", s.account.handleAccount)
	apiGroup.GET("/account/:id", s.account.handleGetAccountById)
	apiGroup.PUT("/account/:id", s.account.handleUpdateAccount)
	apiGroup.DELETE("/account/:id", s.account.handleDeleteAccount)

	apiGroup.Any("/book", s.book.handleBook)
	apiGroup.GET("/book/:id", s.book.handleGetBookById)
	apiGroup.PUT("/book/:id", s.book.handleUpdateBook)
	apiGroup.DELETE("/book/:id", s.book.handleDeleteBook)

	apiGroup.Any("/book/letter", s.book.handleLetter)
	apiGroup.GET("/book/letter/:id", s.book.handleGetLetterById)
	apiGroup.PUT("/book/letter/:id", s.book.handleUpdateLetter)
	apiGroup.DELETE("/book/letter/:id", s.book.handleDeleteLetter)

	apiGroup.Any("/book/version", s.book.handleVersion)
	apiGroup.GET("/book/version/:id", s.book.handleGetVersionById)
	apiGroup.PUT("/book/version/:id", s.book.handleUpdateVersion)
	apiGroup.DELETE("/book/version/:id", s.book.handleDeleteVersion)

	apiGroup.Any("/book/cover", s.book.handleCover)
	apiGroup.GET("/book/cover/:id", s.book.handleGetCoverById)
	apiGroup.PUT("/book/cover/:id", s.book.handleUpdateCover)
	apiGroup.DELETE("/book/cover/:id", s.book.handleDeleteCover)

	apiGroup.Any("/book/publisher", s.book.handlePublisher)
	apiGroup.GET("/book/publisher/:id", s.book.handleGetPublisherById)
	apiGroup.PUT("/book/publisher/:id", s.book.handleUpdatePublisher)
	apiGroup.DELETE("/book/publisher/:id", s.book.handleDeletePublisher)

	apiGroup.Any("/book/bcategory", s.book.handleBCategory)
	apiGroup.GET("/book/bcategory/:id", s.book.handleGetBCategoryById)
	apiGroup.PUT("/book/bcategory/:id", s.book.handleUpdateBCategory)
	apiGroup.DELETE("/book/bcategory/:id", s.book.handleDeleteBCategory)

	apiGroup.Any("/article", s.article.handleArticle)
	apiGroup.GET("/article/:id", s.article.handleGetArticleById)
	apiGroup.PUT("/article/:id", s.article.handleUpdateArticle)
	apiGroup.DELETE("/article/:id", s.article.handleDeleteArticle)

	apiGroup.Any("/article/acategory", s.article.handleACategory)
	apiGroup.GET("/article/acategory/:id", s.article.handleGetACategoryById)
	apiGroup.PUT("/article/acategory/:id", s.article.handleUpdateACategory)
	apiGroup.DELETE("/article/acategory/:id", s.article.handleDeleteACategory)

	apiGroup.Any("/resource", s.resource.handleResource)
	apiGroup.GET("/resource/:id", s.resource.handleGetResourceById)
	apiGroup.PUT("/resource/:id", s.resource.handleUpdateResource)
	apiGroup.DELETE("/resource/:id", s.resource.handleDeleteResource)

	apiGroup.Any("/resource/rcategory", s.resource.handleRCategory)
	apiGroup.GET("/resource/rcategory/:id", s.resource.handleGetRCategoryById)
	apiGroup.PUT("/resource/rcategory/:id", s.resource.handleUpdateRCategory)
	apiGroup.DELETE("/resource/rcategory/:id", s.resource.handleDeleteRCategory)

	apiGroup.Any("/order", s.order.handleOrder)
	apiGroup.GET("/order/:id", s.order.handleGetOrderById)
	apiGroup.PUT("/order/:id", s.order.handleUpdateOrder)
	apiGroup.DELETE("/order/:id", s.order.handleDeleteOrder)
}
