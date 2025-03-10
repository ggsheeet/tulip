package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	allowedOrigin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return func(c echo.Context) error {
		origin := c.Request().Header.Get("Origin")
		referer := c.Request().Header.Get("Referer")

		if referer != "" {
			refererURL, err := url.Parse(referer)
			if err == nil && refererURL.Host != "" {
				pathSegments := strings.Split(strings.Trim(refererURL.Path, "/"), "/")
				lastSegment := pathSegments[len(pathSegments)-1]

				if lastSegment == "cart" || lastSegment == "processed" {
					return next(c)
				}

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

func timeoutMiddleware(timeout time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(c.Request().Context(), timeout)
			defer cancel()

			req := c.Request().Clone(ctx)
			c.SetRequest(req)

			err := next(c)

			if ctx.Err() == context.DeadlineExceeded {
				return echo.NewHTTPError(http.StatusGatewayTimeout, "Request timed out")
			}

			return err
		}
	}
}
