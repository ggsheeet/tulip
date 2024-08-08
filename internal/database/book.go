package database

import "database/sql"

func (s *PostgresDB) createBookTable() error {
	_, err := s.db.Exec(createBookTabQ)

	return err
}

func (s *PostgresDB) CreateBook(book *Book) error {
	_, err := s.db.Query(
		createBookQ,
		&book.Title,
		&book.Author,
		&book.Description,
		&book.CoverURL,
		&book.ISBN,
		&book.Price,
		&book.Stock,
		&book.SalesCount,
		&book.IsActive,
		&book.LetterSizeID,
		&book.VersionID,
		&book.CoverID,
		&book.PublisherID,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) DeleteBook(id string) error {
	_, err := s.db.Exec(deleteBookQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateaBook(*Book) error {
	return nil
}

func (s *PostgresDB) GetBookById(id string) (*Book, error) {
	row := s.db.QueryRow(getBookQ, id)

	var book Book

	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Description,
		&book.CoverURL,
		&book.ISBN,
		&book.Price,
		&book.Stock,
		&book.SalesCount,
		&book.IsActive,
		&book.LetterSizeID,
		&book.VersionID,
		&book.CoverID,
		&book.PublisherID,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &book, nil
}

func (s *PostgresDB) GetBooks() (*[]*Book, error) {
	rows, err := s.db.Query(getBooksQ)
	if err != nil {
		return nil, err
	}
	books := []*Book{}
	for rows.Next() {
		book := new(Book)
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Description,
			&book.CoverURL,
			&book.ISBN,
			&book.Price,
			&book.Stock,
			&book.SalesCount,
			&book.IsActive,
			&book.LetterSizeID,
			&book.VersionID,
			&book.CoverID,
			&book.PublisherID,
			&book.CreatedAt,
			&book.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}
	return &books, nil
}
