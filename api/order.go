package api

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"

	"github.com/ggsheet/tulip/internal/database"
	"github.com/labstack/echo/v4"
)

func (s *OrderHandlers) handleOrder(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetOrders(c)
	case http.MethodPost:
		return s.handleCreateOrder(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
}

func (s *OrderHandlers) handleGetOrders(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	orders, err := s.db.GetOrders(page, limit)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, orders)
}

func (s *OrderHandlers) handleGetUnfulfilledOrders(c echo.Context) error {
	orders, err := s.db.GetUnfulfilledOrders()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, orders)
}

func (s *OrderHandlers) handleGetFulfilledOrders(c echo.Context) error {
	orders, err := s.db.GetFulfilledOrders()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, orders)
}

func (s *OrderHandlers) handleGetOrderById(c echo.Context) error {
	id := c.Param("id")
	order, err := s.db.GetOrderById(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, order)
}

func (s *OrderHandlers) handleGetOrderByIdAdmin(c echo.Context) error {
	id := c.Param("id")
	order, bookOrders, books, err := s.db.GetOrderWithBooksAdmin(id)
	if err != nil {
		fmt.Println("error", err)
		return err
	}

	response := map[string]interface{}{
		"order":      order,
		"bookOrders": bookOrders,
		"books":      books,
	}

	return c.JSON(http.StatusOK, response)
}

func (s *OrderHandlers) handleUpdateOrderStatus(c echo.Context) error {
	id := c.Param("id")
	status := c.FormValue("status")

	if id == "" {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "ID not found"})
	}

	if status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Status is required"})
	}

	// Validate status value
	validStatuses := []string{"processing", "delivered", "returned"}

	if !slices.Contains(validStatuses, status) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid status value"})
	}

	err := s.db.UpdateOrderStatus(id, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update order status"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Estado actualizado correctamente"})
}

func (s *OrderHandlers) handleCreateOrder(c echo.Context) error {
	orderReq := new(database.OrderRequest)

	if err := c.Bind(orderReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order := database.NewOrder(orderReq.Address, orderReq.Total, orderReq.PaymentID, orderReq.IsFulfilled, orderReq.Status, orderReq.AccountID)

	if _, err := s.db.CreateOrder(order); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, order)
}

func (s *OrderHandlers) handleDeleteOrder(c echo.Context) error {
	id := c.Param("id")

	if _, err := s.db.GetOrderById(id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("ID not found, operation unsuccessful: %v", err))
	}

	if err := s.db.DeleteOrder(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func (s *OrderHandlers) handleUpdateOrder(c echo.Context) error {
	id := c.Param("id")
	orderReq := new(database.OrderRequest)

	if err := c.Bind(orderReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order := database.UpdateOrder(orderReq.Address, orderReq.Total, orderReq.PaymentID, orderReq.IsFulfilled, orderReq.Status, orderReq.AccountID)

	if err := s.db.UpdateOrder(id, order); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, order)
}
