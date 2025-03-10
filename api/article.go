package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ggsheet/tulip/internal/database"
	"github.com/labstack/echo/v4"
)

func (s *ArticleHandlers) handleArticle(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetArticles(c)
	case http.MethodPost:
		return s.handleCreateArticle(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *ArticleHandlers) handleGetArticles(c echo.Context) error {
	page := 1
	limit := 10
	category := 0
	order := ""
	articleId := 0

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

	if articleIdParam := c.QueryParam("itemId"); articleIdParam != "" {
		var err error
		articleId, err = strconv.Atoi(articleIdParam)
		if err != nil || articleId <= 0 {
			return c.JSON(http.StatusBadRequest, "Invalid article id")
		}
	}

	articles, err := s.db.GetArticles(page, limit, category, order, articleId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error fetching articles")
	}
	return c.JSON(http.StatusOK, articles)
}

func (s *ArticleHandlers) handleGetArticleById(c echo.Context) error {
	id := c.Param("id")
	article, err := s.db.GetArticleById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, APIError{Error: err.Error()})
	}
	return c.JSON(http.StatusOK, article)
}

func (s *ArticleHandlers) handleCreateArticle(c echo.Context) error {
	artReq := new(database.ArticleRequest)

	if err := c.Bind(artReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	article := database.NewArticle(artReq.Title, artReq.Author, artReq.Excerpt, artReq.Description, artReq.CoverURL, artReq.CategoryID)

	if err := s.db.CreateArticle(article); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, article)
}

func (s *ArticleHandlers) handleDeleteArticle(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetArticleById(id); err != nil {
		return err
	}

	if err := s.db.DeleteArticle(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Article deleted successfully"})
}

func (s *ArticleHandlers) handleUpdateArticle(c echo.Context) error {
	id := c.Param("id")
	artReq := new(database.ArticleRequest)

	if err := c.Bind(artReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	article := database.UpdateArticle(artReq.Title, artReq.Author, artReq.Excerpt, artReq.Description, artReq.CoverURL, artReq.CategoryID)

	if err := s.db.UpdateArticle(id, article); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, article)
}

func (s *ArticleHandlers) handleACategory(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetACategories(c)
	case http.MethodPost:
		return s.handleCreateACategory(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *ArticleHandlers) handleGetACategories(c echo.Context) error {
	aCategories, err := s.db.GetACategories()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, aCategories)
}

func (s *ArticleHandlers) handleGetACategoryById(c echo.Context) error {
	id := c.Param("id")
	aCategory, err := s.db.GetACategoryById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, aCategory)
}

func (s *ArticleHandlers) handleCreateACategory(c echo.Context) error {
	aCategoryReq := new(database.ACategoryRequest)

	if err := c.Bind(aCategoryReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	aCategory := database.NewACategory(aCategoryReq.ArticleCategory, aCategoryReq.IsActive)

	if err := s.db.CreateACategory(aCategory); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, aCategory)
}

func (s *ArticleHandlers) handleDeleteACategory(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetACategoryById(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("ID not found, operation unsuccessful: %v", err))
	}

	if err := s.db.DeleteACategory(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Article category deleted successfully"})
}

func (s *ArticleHandlers) handleUpdateACategory(c echo.Context) error {
	id := c.Param("id")
	aCategoryReq := new(database.ACategoryRequest)

	if err := c.Bind(aCategoryReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	aCategory := database.UpdateACategory(aCategoryReq.ArticleCategory, aCategoryReq.IsActive)

	if err := s.db.UpdateACategory(id, aCategory); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, aCategory)
}
