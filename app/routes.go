package app

import (
	"github.com/labstack/echo/v4"
)

func APPRouter(e *echo.Echo) {
	e.GET("/", handleIndexPage)
	e.GET("/store", handleStorePage)
	e.GET("/book", handleBookPage)
	e.GET("/articles", handleArticlesPage)
	e.GET("/article", handleArticlePage)
	e.GET("/resource", handleResourcePage)
	e.GET("/resources", handleResourcesPage)
}
