package database

import (
	"database/sql"

	"github.com/ggsheet/kerigma/internal/config"
	_ "github.com/lib/pq"
)

func DBConnection() (*PostgresDB, error) {
	connStr := config.GetDatabaseURL()

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{db}, nil
}

func (db *PostgresDB) Init() error {
	errs := []error{
		db.createAccountTable(),
		db.createBookTable(),
		db.createLetterTable(),
		db.createVersionTable(),
		db.createCoverTable(),
		db.createPublisherTable(),
		db.createBCategoryTable(),
		db.createArticleTable(),
		db.createACategoryTable(),
		db.createResourceTable(),
		db.createRCategoryTable(),
		db.createOrderTable(),
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
