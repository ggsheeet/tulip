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

func (s *PostgresDB) GetOrderByIdAdmin(id string) (*Order, error) {
	row := s.db.QueryRow(getOrderByIdAdminQ, id)

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

func (s *PostgresDB) GetOrderWithBooksAdmin(id string) (*Order, []BookOrder, []Book, error) {
	rows, err := s.db.Query(getOrderByIdAdminQ, id)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	var order *Order
	var bookOrders []BookOrder
	var books []Book
	bookMap := make(map[int]Book)
	bookOrderMap := make(map[int]BookOrder)

	for rows.Next() {
		var o Order
		var bo BookOrder
		var b Book
		var bookOrderID, bookID sql.NullInt64
		var bookOrderCreatedAt, bookOrderUpdatedAt, bookCreatedAt, bookUpdatedAt sql.NullTime
		var bookTitle, bookDescription, bookPictureURL sql.NullString
		var bookPrice sql.NullFloat64
		var bookIsActive sql.NullBool

		err := rows.Scan(
			// Order fields
			&o.ID, &o.Address, &o.Total, &o.PaymentID, &o.IsFulfilled, &o.Status, &o.AccountID, &o.CreatedAt, &o.UpdatedAt,
			// BookOrder fields
			&bookOrderID, &bo.Quantity, &bo.BookID, &bo.OrderID, &bookOrderCreatedAt, &bookOrderUpdatedAt,
			// Book fields
			&bookID, &bookTitle, &bookDescription, &bookPrice, &bookPictureURL, &bookIsActive, &bookCreatedAt, &bookUpdatedAt,
		)

		if err != nil {
			return nil, nil, nil, err
		}

		// Set order data (will be the same for all rows)
		if order == nil {
			order = &o
		}

		// Handle BookOrder data
		if bookOrderID.Valid {
			bo.ID = int(bookOrderID.Int64)
			if bookOrderCreatedAt.Valid {
				bo.CreatedAt = bookOrderCreatedAt.Time
			}
			if bookOrderUpdatedAt.Valid {
				bo.UpdatedAt = bookOrderUpdatedAt.Time
			}
			bookOrderMap[bo.ID] = bo
		}

		// Handle Book data
		if bookID.Valid {
			b.ID = int(bookID.Int64)
			if bookTitle.Valid {
				b.Title = bookTitle.String
			}
			if bookDescription.Valid {
				b.Description = bookDescription.String
			}
			if bookPrice.Valid {
				b.Price = bookPrice.Float64
			}
			if bookPictureURL.Valid {
				b.CoverURL = bookPictureURL.String
			}
			if bookIsActive.Valid {
				b.IsActive = bookIsActive.Bool
			}
			if bookCreatedAt.Valid {
				b.CreatedAt = bookCreatedAt.Time
			}
			if bookUpdatedAt.Valid {
				b.UpdatedAt = bookUpdatedAt.Time
			}
			bookMap[b.ID] = b
		}
	}

	// Convert maps to slices
	for _, bo := range bookOrderMap {
		bookOrders = append(bookOrders, bo)
	}
	for _, b := range bookMap {
		books = append(books, b)
	}

	if order == nil {
		return nil, nil, nil, sql.ErrNoRows
	}

	return order, bookOrders, books, nil
}

func (s *PostgresDB) UpdateOrderStatus(id string, status string) error {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(updateOrderStatusQ, status, idInt)
	if err != nil {
		return err
	}

	if status != "processing" {
		_, err = s.db.Exec(fulfillOrderQ, idInt)
		if err != nil {
			return err
		}
	} else {
		_, err = s.db.Exec(unfulfillOrderQ, idInt)
		if err != nil {
			return err
		}
	}

	return err
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

func (s *PostgresDB) GetOrders(page, limit int) (*[]*Order, error) {
	offset := (page - 1) * limit

	query := getOrdersQ + ` ORDER BY o.created_at DESC LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*Order{}
	for rows.Next() {
		order := new(Order)
		var firstName, lastName, email, phone sql.NullString

		err := rows.Scan(
			&order.ID,
			&order.Address,
			&order.Total,
			&order.PaymentID,
			&order.IsFulfilled,
			&order.Status,
			&order.AccountID,
			&firstName,
			&lastName,
			&email,
			&phone,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.RecordCount,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	return &orders, nil
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
