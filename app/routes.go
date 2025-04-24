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
	e.GET("/cart", handleCartPage)
	e.GET("/download", handleResourceDownload)
	e.GET("/processed", handleProcesedPage)
	e.POST("/notification", handlePaymentNotification)
	e.GET("/login", handleLoginPage)
	e.GET("/auth", handleAuthCheck)
	e.POST("/login", handleLoginAuth)
	e.POST("/logout", handleLogoutAuth)
	e.GET("/admin", handleAdminPage)
	e.GET("/sitemap", handleSitemap)

	// Debugging
	e.GET("/test-email", handleTestEmail)
}
