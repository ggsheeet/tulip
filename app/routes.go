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
	e.GET("/admin/table-content", handleAdminTableContent)

	adminGroup := e.Group("/admin")
	adminGroup.Use(adminAuthMiddleware)
	adminGroup.PUT("/book/:action", handleBookAction)
	adminGroup.PUT("/article/:action", handleArticleAction)
	adminGroup.PUT("/resource/:action", handleResourceAction)
	adminGroup.PUT("/bcategory/:action", handleBookCategoryAction)
	adminGroup.PUT("/acategory/:action", handleArticleCategoryAction)
	adminGroup.PUT("/rcategory/:action", handleResourceCategoryAction)
	adminGroup.PUT("/publisher/:action", handlePublisherAction)
	adminGroup.PUT("/version/:action", handleVersionAction)
	adminGroup.PUT("/letter/:action", handleLetterAction)
	adminGroup.PUT("/cover/:action", handleCoverAction)

	// Edit routes
	adminGroup.GET("/edit/book", handleBookEditPage)
	adminGroup.GET("/edit/article", handleArticleEditPage)
	adminGroup.GET("/edit/resource", handleResourceEditPage)
	adminGroup.GET("/edit/bcategory", handleBookCategoryEditPage)
	adminGroup.GET("/edit/acategory", handleArticleCategoryEditPage)
	adminGroup.GET("/edit/rcategory", handleResourceCategoryEditPage)
	adminGroup.GET("/edit/publisher", handlePublisherEditPage)
	adminGroup.GET("/edit/version", handleVersionEditPage)
	adminGroup.GET("/edit/letter", handleLetterEditPage)
	adminGroup.GET("/edit/cover", handleCoverEditPage)

	// Create routes
	adminGroup.GET("/create/book", handleBookCreatePage)
	adminGroup.GET("/create/article", handleArticleCreatePage)
	adminGroup.GET("/create/resource", handleResourceCreatePage)
	adminGroup.GET("/create/bcategory", handleBookCategoryCreatePage)
	adminGroup.GET("/create/acategory", handleArticleCategoryCreatePage)
	adminGroup.GET("/create/rcategory", handleResourceCategoryCreatePage)
	adminGroup.GET("/create/publisher", handlePublisherCreatePage)
	adminGroup.GET("/create/version", handleVersionCreatePage)
	adminGroup.GET("/create/letter", handleLetterCreatePage)
	adminGroup.GET("/create/cover", handleCoverCreatePage)

	// Update routes
	adminGroup.POST("/book/update", handleBookUpdate)
	adminGroup.POST("/article/update", handleArticleUpdate)
	adminGroup.POST("/resource/update", handleResourceUpdate)
	adminGroup.POST("/bcategory/update", handleBookCategoryUpdate)
	adminGroup.POST("/acategory/update", handleArticleCategoryUpdate)
	adminGroup.POST("/rcategory/update", handleResourceCategoryUpdate)
	adminGroup.POST("/publisher/update", handlePublisherUpdate)
	adminGroup.POST("/version/update", handleVersionUpdate)
	adminGroup.POST("/letter/update", handleLetterUpdate)
	adminGroup.POST("/cover/update", handleCoverUpdate)
	adminGroup.POST("/order/status", handleOrderStatusUpdate)

	// Create POST routes
	adminGroup.POST("/book/create", handleBookCreate)
	adminGroup.POST("/article/create", handleArticleCreate)
	adminGroup.POST("/resource/create", handleResourceCreate)
	adminGroup.POST("/bcategory/create", handleBookCategoryCreate)
	adminGroup.POST("/acategory/create", handleArticleCategoryCreate)
	adminGroup.POST("/rcategory/create", handleResourceCategoryCreate)
	adminGroup.POST("/publisher/create", handlePublisherCreate)
	adminGroup.POST("/version/create", handleVersionCreate)
	adminGroup.POST("/letter/create", handleLetterCreate)
	adminGroup.POST("/cover/create", handleCoverCreate)

	// Order routes
	adminGroup.GET("/order", handleOrderPage)
	adminGroup.GET("/view/order", handleOrderPage)

	e.GET("/sitemap", handleSitemap)

	// Debugging
	e.GET("/test-email", handleTestEmail)
}
