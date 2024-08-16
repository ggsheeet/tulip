package app

import (
	"github.com/labstack/echo/v4"
)

func APPRouter(e *echo.Echo) {
	e.GET("/", handleIndexPage)
}
