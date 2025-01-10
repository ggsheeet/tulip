package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ggsheet/kerigma/internal/database"
	"github.com/labstack/echo/v4"
)

func (s *BookHandlers) handleBook(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetBooks(c)
	case http.MethodPost:
		return s.handleCreateBook(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *BookHandlers) handleGetBooks(c echo.Context) error {
	page := 1
	limit := 10
	category := 0
	order := ""
	bookId := 0
	bookIds := ""

	if pageParam := c.QueryParam("page"); pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil || page <= 0 {
			return c.JSON(http.StatusBadRequest, "Invalid page number")
		}
	}

	if limitParam := c.QueryParam("limit"); limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			return c.JSON(http.StatusBadRequest, "Invalid limit number")
		}
	}

	if categoryParam := c.QueryParam("category"); categoryParam != "" {
		var err error
		category, err = strconv.Atoi(categoryParam)
		if err != nil || category <= 0 {
			return c.JSON(http.StatusBadRequest, "Invalid cateogry id")
		}
	}

	if orderParam := c.QueryParam("order"); orderParam != "" {
		order = orderParam
	}

	if bookIdParam := c.QueryParam("itemId"); bookIdParam != "" {
		var err error
		bookId, err = strconv.Atoi(bookIdParam)
		if err != nil || bookId <= 0 {
			return c.JSON(http.StatusBadRequest, "Invalid book id")
		}
	}

	if bookIdsParam := c.QueryParam("itemIds"); bookIdsParam != "" {
		limit = 999
		bookIds = bookIdsParam
	}

	books, err := s.db.GetBooks(page, limit, category, order, bookId, bookIds)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error fetching books")
	}
	return c.JSON(http.StatusOK, books)
}

func (s *BookHandlers) handleGetBookById(c echo.Context) error {
	id := c.Param("id")
	book, err := s.db.GetBookById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, book)
}

func (s *BookHandlers) handleCreateBook(c echo.Context) error {
	bookReq := new(database.BookRequest)

	if err := c.Bind(bookReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	book := database.NewBook(bookReq.Title, bookReq.Author, bookReq.Description, bookReq.CoverURL, bookReq.ISBN, bookReq.Price, bookReq.Stock, bookReq.SalesCount, bookReq.IsActive, bookReq.LetterID, bookReq.VersionID, bookReq.CoverID, bookReq.PublisherID, bookReq.CategoryID)

	if err := s.db.CreateBook(book); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, book)
}

func (s *BookHandlers) handleDeleteBook(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetBookById(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("ID not found, operation unsuccessful: %v", err))
	}

	if err := s.db.DeleteBook(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Book deleted successfully"})
}

func (s *BookHandlers) handleUpdateBook(c echo.Context) error {
	id := c.Param("id")
	bookReq := new(database.BookRequest)

	if err := c.Bind(bookReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	book := database.UpdateBook(bookReq.Title, bookReq.Author, bookReq.Description, bookReq.CoverURL, bookReq.ISBN, bookReq.Price, bookReq.Stock, bookReq.SalesCount, bookReq.IsActive, bookReq.LetterID, bookReq.VersionID, bookReq.CoverID, bookReq.PublisherID, bookReq.CategoryID)

	if err := s.db.UpdateBook(id, book); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, book)
}

func (s *BookHandlers) handleLetter(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetLetters(c)
	case http.MethodPost:
		return s.handleCreateLetter(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *BookHandlers) handleGetLetters(c echo.Context) error {
	letters, err := s.db.GetLetters()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, letters)
}

func (s *BookHandlers) handleGetLetterById(c echo.Context) error {
	id := c.Param("id")
	letter, err := s.db.GetLetterById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, letter)
}

func (s *BookHandlers) handleCreateLetter(c echo.Context) error {
	letterReq := new(database.LetterRequest)

	if err := c.Bind(letterReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	letter := database.NewLetter(letterReq.LetterType)

	if err := s.db.CreateLetter(letter); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, letter)
}

func (s *BookHandlers) handleDeleteLetter(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetLetterById(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("ID not found, operation unsuccessful: %v", err))
	}

	if err := s.db.DeleteLetter(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Letter deleted successfully"})
}

func (s *BookHandlers) handleUpdateLetter(c echo.Context) error {
	id := c.Param("id")
	letterReq := new(database.LetterRequest)

	if err := c.Bind(letterReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	letter := database.NewLetter(letterReq.LetterType)

	if err := s.db.UpdateLetter(id, letter); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, letter)
}

func (s *BookHandlers) handleVersion(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetVersions(c)
	case http.MethodPost:
		return s.handleCreateVersion(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *BookHandlers) handleGetVersions(c echo.Context) error {
	versions, err := s.db.GetVersions()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, versions)
}

func (s *BookHandlers) handleGetVersionById(c echo.Context) error {
	id := c.Param("id")
	version, err := s.db.GetVersionById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, version)
}

func (s *BookHandlers) handleCreateVersion(c echo.Context) error {
	versionReq := new(database.VersionRequest)

	if err := c.Bind(versionReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	version := database.NewVersion(versionReq.BibleVersion)

	if err := s.db.CreateVersion(version); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, version)
}

func (s *BookHandlers) handleDeleteVersion(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetVersionById(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("ID not found, operation unsuccessful: %v", err))
	}

	if err := s.db.DeleteVersion(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Version deleted successfully"})
}

func (s *BookHandlers) handleUpdateVersion(c echo.Context) error {
	id := c.Param("id")
	versionReq := new(database.VersionRequest)

	if err := c.Bind(versionReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	version := database.UpdateVersion(versionReq.BibleVersion)

	if err := s.db.UpdateVersion(id, version); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, version)
}

func (s *BookHandlers) handleCover(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetCovers(c)
	case http.MethodPost:
		return s.handleCreateCover(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *BookHandlers) handleGetCovers(c echo.Context) error {
	covers, err := s.db.GetCovers()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, covers)
}

func (s *BookHandlers) handleGetCoverById(c echo.Context) error {
	id := c.Param("id")
	cover, err := s.db.GetCoverById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cover)
}

func (s *BookHandlers) handleCreateCover(c echo.Context) error {
	coverReq := new(database.CoverRequest)

	if err := c.Bind(coverReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	cover := database.NewCover(coverReq.CoverType)

	if err := s.db.CreateCover(cover); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cover)
}

func (s *BookHandlers) handleDeleteCover(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetCoverById(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("ID not found, operation unsuccessful: %v", err))
	}

	if err := s.db.DeleteCover(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Cover deleted successfully"})
}

func (s *BookHandlers) handleUpdateCover(c echo.Context) error {
	id := c.Param("id")
	coverReq := new(database.CoverRequest)

	if err := c.Bind(coverReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	cover := database.UpdateCover(coverReq.CoverType)

	if err := s.db.UpdateCover(id, cover); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cover)
}

func (s *BookHandlers) handlePublisher(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetPublishers(c)
	case http.MethodPost:
		return s.handleCreatePublisher(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *BookHandlers) handleGetPublishers(c echo.Context) error {
	publishers, err := s.db.GetPublishers()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, publishers)
}

func (s *BookHandlers) handleGetPublisherById(c echo.Context) error {
	id := c.Param("id")
	publisher, err := s.db.GetPublisherById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, publisher)
}

func (s *BookHandlers) handleCreatePublisher(c echo.Context) error {
	publisherReq := new(database.PublisherRequest)

	if err := c.Bind(publisherReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	publisher := database.NewPublisher(publisherReq.PublisherName)

	if err := s.db.CreatePublisher(publisher); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, publisher)
}

func (s *BookHandlers) handleDeletePublisher(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetPublisherById(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("ID not found, operation unsuccessful: %v", err))
	}

	if err := s.db.DeletePublisher(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Publisher deleted successfully"})
}

func (s *BookHandlers) handleUpdatePublisher(c echo.Context) error {
	id := c.Param("id")
	publisherReq := new(database.PublisherRequest)

	if err := c.Bind(publisherReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	publisher := database.UpdatePublisher(publisherReq.PublisherName)

	if err := s.db.UpdatePublisher(id, publisher); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, publisher)
}

func (s *BookHandlers) handleBCategory(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetBCategories(c)
	case http.MethodPost:
		return s.handleCreateBCategory(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *BookHandlers) handleGetBCategories(c echo.Context) error {
	bCategories, err := s.db.GetBCategories()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, bCategories)
}

func (s *BookHandlers) handleGetBCategoryById(c echo.Context) error {
	id := c.Param("id")
	bCategory, err := s.db.GetBCategoryById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, bCategory)
}

func (s *BookHandlers) handleCreateBCategory(c echo.Context) error {
	bCategoryReq := new(database.BCategoryRequest)

	if err := c.Bind(bCategoryReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	bCategory := database.NewBCategory(bCategoryReq.BookCategory)

	if err := s.db.CreateBCategory(bCategory); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, bCategory)
}

func (s *BookHandlers) handleDeleteBCategory(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetBCategoryById(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("ID not found, operation unsuccessful: %v", err))
	}

	if err := s.db.DeleteBCategory(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "BCategory deleted successfully"})
}

func (s *BookHandlers) handleUpdateBCategory(c echo.Context) error {
	id := c.Param("id")
	bCategoryReq := new(database.BCategoryRequest)

	if err := c.Bind(bCategoryReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	bCategory := database.UpdateBCategory(bCategoryReq.BookCategory)

	if err := s.db.UpdateBCategory(id, bCategory); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, bCategory)
}
