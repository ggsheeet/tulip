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

func (s *PostgresDB) CreateOrder(order *Order) error {
	_, err := s.db.Query(
		createOrderQ,
		&order.FirstName,
		&order.LastName,
		&order.Address,
		&order.Quantity,
		&order.Total,
		&order.BookID,
		&order.AccountID,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
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
		&order.FirstName,
		&order.LastName,
		&order.Address,
		&order.Quantity,
		&order.Total,
		&order.BookID,
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
		&order.FirstName,
		&order.LastName,
		&order.Address,
		&order.Quantity,
		&order.Total,
		&order.BookID,
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

func (s *PostgresDB) GetOrders() (*[]*Order, error) {
	rows, err := s.db.Query(getOrdersQ)
	if err != nil {
		return nil, err
	}
	orders := []*Order{}
	for rows.Next() {
		order := new(Order)
		err := rows.Scan(
			&order.ID,
			&order.FirstName,
			&order.LastName,
			&order.Address,
			&order.Quantity,
			&order.Total,
			&order.BookID,
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
