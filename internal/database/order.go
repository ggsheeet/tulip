package database

import (
	"database/sql"
	"strconv"
	"time"
)

func (s *PostgresDB) createOrderTable() error {
	_, err := s.db.Exec(createOrderTabQ)

	return err
}

func (s *PostgresDB) CreateOrder(order *Order) (int, error) {
	err := s.db.QueryRow(
		createOrderQ,
		&order.Address,
		&order.Total,
		&order.PaymentID,
		&order.IsFulfilled,
		&order.Status,
		&order.AccountID,
		&order.CreatedAt,
		&order.UpdatedAt,
	).Scan(&order.ID)

	if err != nil {
		return 0, err
	}

	return order.ID, nil
}

func (s *PostgresDB) DeleteOrder(id string) error {
	_, err := s.db.Exec(deleteOrderQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateOrder(id string, order *Order) error {
	orderId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}

	_, err := s.db.Query(
		updateOrderQ,
		orderId,
		&order.Address,
		&order.Total,
		&order.PaymentID,
		&order.IsFulfilled,
		&order.Status,
		&order.AccountID,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetOrderById(id string) (*Order, error) {
	row := s.db.QueryRow(getOrderQ, id)

	var order Order

	err := row.Scan(
		&order.ID,
		&order.Address,
		&order.Total,
		&order.PaymentID,
		&order.IsFulfilled,
		&order.Status,
		&order.AccountID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &order, nil
}

func (s *PostgresDB) GetOrderByPaymentId(paymentId int) (int, error) {
	row := s.db.QueryRow(getOrderByPaymentIdQ, paymentId)

	var orderID int
	err := row.Scan(&orderID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return orderID, nil
}

func (s *PostgresDB) GetUnfulfilledOrders() (*[]*Order, error) {
	rows, err := s.db.Query(getUnfulfilledOrdersQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*Order{}
	for rows.Next() {
		order := new(Order)
		err := rows.Scan(
			&order.ID,
			&order.Address,
			&order.Total,
			&order.PaymentID,
			&order.IsFulfilled,
			&order.Status,
			&order.AccountID,
			&order.AccountID,
			&order.CreatedAt,
			&order.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	return &orders, nil
}

func (s *PostgresDB) GetFulfilledOrders() (*[]*Order, error) {
	rows, err := s.db.Query(getFulfilledOrdersQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*Order{}
	for rows.Next() {
		order := new(Order)
		err := rows.Scan(
			&order.ID,
			&order.Address,
			&order.Total,
			&order.PaymentID,
			&order.IsFulfilled,
			&order.Status,
			&order.AccountID,
			&order.AccountID,
			&order.CreatedAt,
			&order.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	return &orders, nil
}

func (s *PostgresDB) createBookOrderTable() error {
	_, err := s.db.Exec(createBookOrdersTabQ)

	return err
}

func (s *PostgresDB) CreateBookOrder(bookQuantity int, bookId int, orderId int) error {
	bookOrder := NewBookOrder(bookQuantity, bookId, orderId)

	_, err := s.db.Exec(
		createBookOrderQ,
		&bookOrder.Quantity,
		&bookOrder.BookID,
		&bookOrder.OrderID,
		&bookOrder.CreatedAt,
		&bookOrder.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
