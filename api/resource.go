package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ggsheet/kerigma/internal/database"
	"github.com/labstack/echo/v4"
)

func (s *ResourceHandlers) handleResource(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetResources(c)
	case http.MethodPost:
		return s.handleCreateResource(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *ResourceHandlers) handleGetResources(c echo.Context) error {
	page := 1
	limit := 10
	category := 0
	order := ""
	resourceId := 0

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

	if resourceIdParam := c.QueryParam("itemId"); resourceIdParam != "" {
		var err error
		resourceId, err = strconv.Atoi(resourceIdParam)
		if err != nil || resourceId <= 0 {
			return c.JSON(http.StatusBadRequest, "Invalid resource id")
		}
	}

	resources, err := s.db.GetResources(page, limit, category, order, resourceId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error fetching resources")
	}
	return c.JSON(http.StatusOK, resources)
}

func (s *ResourceHandlers) handleGetResourceById(c echo.Context) error {
	id := c.Param("id")
	resource, err := s.db.GetResourceById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resource)
}

func (s *ResourceHandlers) handleCreateResource(c echo.Context) error {
	resReq := new(database.ResourceRequest)

	if err := c.Bind(resReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resource := database.NewResource(resReq.Title, resReq.Author, resReq.Description, resReq.CoverURL, resReq.ResourceURL, resReq.CategoryID)

	if err := s.db.CreateResource(resource); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resource)
}

func (s *ResourceHandlers) handleDeleteResource(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetResourceById(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("ID not found, operation unsuccessful: %v", err))
	}

	if err := s.db.DeleteResource(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Resource deleted successfully"})
}

func (s *ResourceHandlers) handleUpdateResource(c echo.Context) error {
	id := c.Param("id")
	resReq := new(database.ResourceRequest)

	if err := c.Bind(resReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	resource := database.UpdateResource(resReq.Title, resReq.Author, resReq.Description, resReq.CoverURL, resReq.ResourceURL, resReq.CategoryID)

	if err := s.db.UpdateResource(id, resource); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resource)
}

func (s *ResourceHandlers) handleRCategory(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetRCategories(c)
	case http.MethodPost:
		return s.handleCreateRCategory(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *ResourceHandlers) handleGetRCategories(c echo.Context) error {
	aCategories, err := s.db.GetRCategories()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, aCategories)
}

func (s *ResourceHandlers) handleGetRCategoryById(c echo.Context) error {
	id := c.Param("id")
	rCategory, err := s.db.GetRCategoryById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, rCategory)
}

func (s *ResourceHandlers) handleCreateRCategory(c echo.Context) error {
	rCategoryReq := new(database.RCategoryRequest)

	if err := c.Bind(rCategoryReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	rCategory := database.NewRCategory(rCategoryReq.ResourceCategory)

	if err := s.db.CreateRCategory(rCategory); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, rCategory)
}

func (s *ResourceHandlers) handleDeleteRCategory(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetRCategoryById(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("ID not found, operation unsuccessful: %v", err))
	}

	if err := s.db.DeleteRCategory(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Resource category deleted successfully"})
}

func (s *ResourceHandlers) handleUpdateRCategory(c echo.Context) error {
	id := c.Param("id")
	rCategoryReq := new(database.RCategoryRequest)

	if err := c.Bind(rCategoryReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	rCategory := database.UpdateRCategory(rCategoryReq.ResourceCategory)

	if err := s.db.UpdateRCategory(id, rCategory); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, rCategory)
}
