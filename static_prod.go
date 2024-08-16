//go:build !dev
// +build !dev

package main

import (
	"embed"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:embed public
var publicFS embed.FS

func public() http.Handler {
	return http.FileServer(http.FS(publicFS))
}

// Wrap the http.Handler for Echo
func echoWrapHandler(h http.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}
