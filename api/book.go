package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ggsheet/kerigma/internal/database"
)

func (s *BookHandlers) handleBook(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetBooks(w)
	case "POST":
		return s.handleCreateBook(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *BookHandlers) handleGetBooks(w http.ResponseWriter) error {
	books, err := s.bookInterface.GetBooks()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, books)
}

func (s *BookHandlers) handleGetBookById(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	book, err := s.bookInterface.GetBookById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, book)
}

func (s *BookHandlers) handleCreateBook(w http.ResponseWriter, r *http.Request) error {
	cbr := new(database.CreateBookRequest)

	if err := json.NewDecoder(r.Body).Decode(cbr); err != nil {
		return err
	}

	book := database.NewBook(cbr.Title, cbr.Author, cbr.Description, cbr.CoverURL, cbr.ISBN, cbr.Price, cbr.Stock, cbr.SalesCount, cbr.IsActive, cbr.LetterSizeID, cbr.VersionID, cbr.CoverID, cbr.CategoryID, cbr.PublisherID)

	if err := s.bookInterface.CreateBook(book); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, book)
}

func (s *BookHandlers) handleDeleteBook(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")

	if _, err := s.bookInterface.GetBookById(id); err != nil {
		return fmt.Errorf("ID not found, operation unsuccessful")
	} else {
		err := s.bookInterface.DeleteBook(id)

		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, "Book deleted successfully")
	}
}

func (s *BookHandlers) handlUpdateBook(w http.ResponseWriter, r *http.Request) error {

	return nil
}
