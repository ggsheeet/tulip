//+build dev
//go:build dev
// +build dev

package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func public() http.Handler {
	return http.StripPrefix("/public/", http.FileServer(http.FS(os.DirFS("public"))))
}

// Wrap the http.Handler for Echo
func echoWrapHandler(h http.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	}
}
