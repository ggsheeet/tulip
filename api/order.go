package api

import (
	"fmt"
	"net/http"

	"github.com/ggsheet/tulip/internal/database"
	"github.com/labstack/echo/v4"
)

func (s *OrderHandlers) handleOrder(c echo.Context) error {
	switch c.Request().Method {
	case http.MethodGet:
		return s.handleGetUnfulfilledOrders(c)
	case http.MethodPost:
		return s.handleCreateOrder(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed %s", c.Request().Method))
	}
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
