package app

import (
	"github.com/ggsheet/kerigma/template/layout"
	"github.com/labstack/echo/v4"
)

func handleIndexPage(c echo.Context) error {
	return Render(c, layout.Index())
}
