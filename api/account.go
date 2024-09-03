package api

import (
	"fmt"
	"net/http"

	"github.com/ggsheet/kerigma/internal/database"
	"github.com/labstack/echo/v4"
)

func (s *AccountHandlers) handleAccount(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetAccounts(c)
	case http.MethodPost:
		return s.handleCreateAccount(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *AccountHandlers) handleGetAccounts(c echo.Context) error {
	accounts, err := s.db.GetAccounts()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, accounts)
}

func (s *AccountHandlers) handleGetAccountById(c echo.Context) error {
	id := c.Param("id")
	account, err := s.db.GetAccountById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIError{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, account)
}

func (s *AccountHandlers) handleCreateAccount(c echo.Context) error {
	accReq := new(database.AccountRequest)

	if err := c.Bind(accReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	account := database.NewAccount(accReq.FirstName, accReq.LastName, accReq.Email)

	if err := s.db.CreateAccount(account); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, account)
}

func (s *AccountHandlers) handleDeleteAccount(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetAccountById(id); err != nil {
		return err
	}

	if err := s.db.DeleteAccount(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (s *AccountHandlers) handleUpdateAccount(c echo.Context) error {
	id := c.Param("id")
	accReq := new(database.AccountRequest)

	if err := c.Bind(accReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	account := database.UpdateAccount(accReq.FirstName, accReq.LastName, accReq.Email)

	if err := s.db.UpdateAccount(id, account); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, account)
}
