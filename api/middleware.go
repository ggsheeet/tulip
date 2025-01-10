package api

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
)

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	environment := os.Getenv("ENVIRONMENT")
	var allowedOrigin string
	token := os.Getenv("AUTH_TOKEN")

	if environment == "development" || environment == "docker" {
		allowedOrigin = os.Getenv("AUTH_ORIGIN_DEV")
	} else if environment == "production" {
		allowedOrigin = os.Getenv("AUTH_ORIGIN_PROD")
	}

	return func(c echo.Context) error {
		origin := c.Request().Header.Get("Origin")
		referer := c.Request().Header.Get("Referer")

		if referer != "" {
			if referer == allowedOrigin+"/cart" {
				return next(c)
			}
			refererURL, err := url.Parse(referer)
			if err == nil && refererURL.Host != "" {
				referer = fmt.Sprintf("%s://%s", refererURL.Scheme, refererURL.Host)
			}
		}

		if origin != allowedOrigin && referer != allowedOrigin {
			return echo.NewHTTPError(http.StatusForbidden, "Forbidden: Invalid Origin or Referer")
		}

		// CORS headers
		c.Response().Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		// Handle preflight requests
		if c.Request().Method == http.MethodOptions {
			return c.NoContent(http.StatusNoContent)
		}

		// Token-based authentication
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "Bearer "+token {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: Invalid Token")
		}

		return next(c)
	}
}
