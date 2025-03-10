package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ggsheet/tulip/internal/database"
	"github.com/google/uuid"
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
		fmt.Printf("Error getting accounts: %v\n", err)
		return err
	}
	return c.JSON(http.StatusOK, accounts)
}

func (s *AccountHandlers) handleGetAccountById(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	account, err := s.db.GetAccountById(id)
	if err != nil {
		fmt.Printf("Error finding account: %v\n", err)
		return c.JSON(http.StatusInternalServerError, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, account)
}

func (s *AccountHandlers) handleGetAccountByEmail(c echo.Context) error {
	type EmailData struct {
		Email string `json:"email"`
	}
	var emailData EmailData

	if err := c.Bind(&emailData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	account, err := s.db.GetAccountByEmail(emailData.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.NoContent(http.StatusNoContent)
		}
		fmt.Printf("Error finding account: %v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error")
	}

	return c.JSON(http.StatusOK, account)
}

func (s *AccountHandlers) handleCreateAccount(c echo.Context) error {
	accReq := new(database.AccountRequest)

	if err := c.Bind(accReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	account := database.NewAccount(accReq.FirstName, accReq.LastName, accReq.Email, accReq.Phone)

	if _, err := s.db.CreateAccount(account); err != nil {
		fmt.Printf("Error creating account: %v\n", err)
		return err
	}

	return c.JSON(http.StatusOK, account)
}

func (s *AccountHandlers) handleDeleteAccount(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	if _, err := s.db.GetAccountById(id); err != nil {
		fmt.Printf("Error deleting account: %v\n", err)
		return err
	}

	if err := s.db.DeleteAccount(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (s *AccountHandlers) handleUpdateAccount(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}
	accReq := new(database.AccountRequest)

	if err := c.Bind(accReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	account := database.UpdateAccount(accReq.FirstName, accReq.LastName, accReq.Email, accReq.Phone)

	if err := s.db.UpdateAccount(id, account); err != nil {
		fmt.Printf("Error updating account: %v\n", err)
		return err
	}
	return c.JSON(http.StatusOK, account)
}
