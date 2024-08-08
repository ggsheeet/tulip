package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ggsheet/kerigma/internal/database"
)

func (s *AccountHandlers) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccounts(w)
	case "POST":
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *AccountHandlers) handleGetAccounts(w http.ResponseWriter) error {
	accounts, err := s.accountInterface.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *AccountHandlers) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	account, err := s.accountInterface.GetAccountById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (s *AccountHandlers) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	car := new(database.CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(car); err != nil {
		return err
	}

	account := database.NewAccount(car.FirstName, car.LastName, car.Email)

	if err := s.accountInterface.CreateAccount(account); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (s *AccountHandlers) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	if _, err := s.accountInterface.GetAccountById(id); err != nil {
		return fmt.Errorf("ID not found, operation unsuccessful")
	} else {
		err := s.accountInterface.DeleteAccount(id)

		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, "User deleted successfully")
	}
}

func (s *AccountHandlers) handlUpdateAccount(w http.ResponseWriter, r *http.Request) error {

	return nil
}
