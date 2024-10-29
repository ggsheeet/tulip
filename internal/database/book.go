package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

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
		&book.LetterID,
		&book.VersionID,
		&book.CoverID,
		&book.PublisherID,
		&book.CategoryID,
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

func (s *PostgresDB) UpdateBook(id string, book *Book) error {
	bookId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}
	_, err := s.db.Query(
		updateBookQ,
		bookId,
		&book.Title,
		&book.Author,
		&book.Description,
		&book.CoverURL,
		&book.ISBN,
		&book.Price,
		&book.Stock,
		&book.SalesCount,
		&book.IsActive,
		&book.LetterID,
		&book.VersionID,
		&book.CoverID,
		&book.PublisherID,
		&book.CategoryID,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

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
		&book.LetterID,
		&book.LetterType,
		&book.VersionID,
		&book.BibleVersion,
		&book.CoverID,
		&book.CoverType,
		&book.PublisherID,
		&book.PublisherName,
		&book.CategoryID,
		&book.BookCategory,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Book with id '%v' does not exist", id)
		}
		return nil, err
	}

	return &book, nil
}

func (s *PostgresDB) GetBooks(page int, limit int, category int, order string, bookId int) (*[]*Book, error) {
	offset := (page - 1) * limit

	query := getBooksQ

	whereClause := " WHERE b.is_active = true"

	if category != 0 {
		whereClause += " AND b.category_id = $3"
	}

	if bookId != 0 {
		whereClause += " AND b.id != $4"
	}

	query += whereClause

	switch order {
	case "expensive":
		query += " ORDER BY b.price DESC"
	case "cheap":
		query += " ORDER BY b.price ASC"
	case "selling":
		query += " ORDER BY b.sales_count DESC"
	default:
		query += " ORDER BY b.created_at DESC, b.id DESC"
	}

	query += " LIMIT $1 OFFSET $2"

	args := []interface{}{limit, offset}
	if category != 0 {
		args = append(args, category)
	}
	if bookId != 0 {
		args = append(args, bookId)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			&book.LetterID,
			&book.LetterType,
			&book.VersionID,
			&book.BibleVersion,
			&book.CoverID,
			&book.CoverType,
			&book.PublisherID,
			&book.PublisherName,
			&book.CategoryID,
			&book.BookCategory,
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

func (s *PostgresDB) createLetterTable() error {
	_, err := s.db.Exec(createLetterTabQ)

	return err
}

func (s *PostgresDB) CreateLetter(letter *Letter) error {
	_, err := s.db.Query(
		createLetterQ,
		&letter.LetterType,
		&letter.CreatedAt,
		&letter.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) DeleteLetter(id string) error {
	_, err := s.db.Exec(deleteLetterQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateLetter(id string, letter *Letter) error {
	letterId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}

	_, err := s.db.Query(
		updateLetterQ,
		letterId,
		&letter.LetterType,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetLetterById(id string) (*Letter, error) {
	row := s.db.QueryRow(getLetterQ, id)

	var letter Letter

	err := row.Scan(
		&letter.ID,
		&letter.LetterType,
		&letter.CreatedAt,
		&letter.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &letter, nil
}

func (s *PostgresDB) GetLetters() (*[]*Letter, error) {
	rows, err := s.db.Query(getLettersQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	letters := []*Letter{}
	for rows.Next() {
		letter := new(Letter)
		err := rows.Scan(
			&letter.ID,
			&letter.LetterType,
			&letter.CreatedAt,
			&letter.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		letters = append(letters, letter)
	}
	return &letters, nil
}

func (s *PostgresDB) createVersionTable() error {
	_, err := s.db.Exec(createVersionTabQ)

	return err
}

func (s *PostgresDB) CreateVersion(version *Version) error {
	_, err := s.db.Query(
		createVersionQ,
		&version.BibleVersion,
		&version.CreatedAt,
		&version.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) DeleteVersion(id string) error {
	_, err := s.db.Exec(deleteVersionQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateVersion(id string, version *Version) error {
	versionId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}

	_, err := s.db.Query(
		updateVersionQ,
		versionId,
		&version.BibleVersion,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetVersionById(id string) (*Version, error) {
	row := s.db.QueryRow(getVersionQ, id)

	var version Version

	err := row.Scan(
		&version.ID,
		&version.BibleVersion,
		&version.CreatedAt,
		&version.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &version, nil
}

func (s *PostgresDB) GetVersions() (*[]*Version, error) {
	rows, err := s.db.Query(getVersionsQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := []*Version{}
	for rows.Next() {
		version := new(Version)
		err := rows.Scan(
			&version.ID,
			&version.BibleVersion,
			&version.CreatedAt,
			&version.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		versions = append(versions, version)
	}
	return &versions, nil
}

func (s *PostgresDB) createCoverTable() error {
	_, err := s.db.Exec(createCoverTabQ)

	return err
}

func (s *PostgresDB) CreateCover(cover *Cover) error {
	_, err := s.db.Query(
		createCoverQ,
		&cover.CoverType,
		&cover.CreatedAt,
		&cover.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) DeleteCover(id string) error {
	_, err := s.db.Exec(deleteCoverQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateCover(id string, cover *Cover) error {
	coverId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}

	_, err := s.db.Query(
		updateCoverQ,
		coverId,
		&cover.CoverType,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetCoverById(id string) (*Cover, error) {
	row := s.db.QueryRow(getCoverQ, id)

	var cover Cover

	err := row.Scan(
		&cover.ID,
		&cover.CoverType,
		&cover.CreatedAt,
		&cover.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &cover, nil
}

func (s *PostgresDB) GetCovers() (*[]*Cover, error) {
	rows, err := s.db.Query(getCoversQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	covers := []*Cover{}
	for rows.Next() {
		cover := new(Cover)
		err := rows.Scan(
			&cover.ID,
			&cover.CoverType,
			&cover.CreatedAt,
			&cover.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		covers = append(covers, cover)
	}
	return &covers, nil
}

func (s *PostgresDB) createPublisherTable() error {
	_, err := s.db.Exec(createPublisherTabQ)

	return err
}

func (s *PostgresDB) CreatePublisher(publisher *Publisher) error {
	_, err := s.db.Query(
		createPublisherQ,
		&publisher.PublisherName,
		&publisher.CreatedAt,
		&publisher.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) DeletePublisher(id string) error {
	_, err := s.db.Exec(deletePublisherQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdatePublisher(id string, publisher *Publisher) error {
	publisherId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}

	_, err := s.db.Query(
		updatePublisherQ,
		publisherId,
		&publisher.PublisherName,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetPublisherById(id string) (*Publisher, error) {
	row := s.db.QueryRow(getPublisherQ, id)

	var publisher Publisher

	err := row.Scan(
		&publisher.ID,
		&publisher.PublisherName,
		&publisher.CreatedAt,
		&publisher.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &publisher, nil
}

func (s *PostgresDB) GetPublishers() (*[]*Publisher, error) {
	rows, err := s.db.Query(getPublishersQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	publishers := []*Publisher{}
	for rows.Next() {
		publisher := new(Publisher)
		err := rows.Scan(
			&publisher.ID,
			&publisher.PublisherName,
			&publisher.CreatedAt,
			&publisher.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		publishers = append(publishers, publisher)
	}
	return &publishers, nil
}

func (s *PostgresDB) createBCategoryTable() error {
	_, err := s.db.Exec(createBCategoryTabQ)

	return err
}

func (s *PostgresDB) CreateBCategory(bCategory *BCategory) error {
	_, err := s.db.Query(
		createBCategoryQ,
		&bCategory.BookCategory,
		&bCategory.CreatedAt,
		&bCategory.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) DeleteBCategory(id string) error {
	_, err := s.db.Exec(deleteBCategoryQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateBCategory(id string, bCategory *BCategory) error {
	bCategoryId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}

	_, err := s.db.Query(
		updateBCategoryQ,
		bCategoryId,
		&bCategory.BookCategory,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetBCategoryById(id string) (*BCategory, error) {
	row := s.db.QueryRow(getBCategoryQ, id)

	var bCategory BCategory

	err := row.Scan(
		&bCategory.ID,
		&bCategory.BookCategory,
		&bCategory.CreatedAt,
		&bCategory.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &bCategory, nil
}

func (s *PostgresDB) GetBCategories() (*[]*BCategory, error) {
	rows, err := s.db.Query(getBCategoriesQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bCategories := []*BCategory{}
	for rows.Next() {
		bCategory := new(BCategory)
		err := rows.Scan(
			&bCategory.ID,
			&bCategory.BookCategory,
			&bCategory.CreatedAt,
			&bCategory.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		bCategories = append(bCategories, bCategory)
	}
	return &bCategories, nil
}
