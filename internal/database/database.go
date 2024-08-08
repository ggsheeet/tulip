package database

import (
	"database/sql"
	"log"

	"github.com/ggsheet/kerigma/internal/config"
	_ "github.com/lib/pq"
)

func DBConnection() (*PostgresDB, error) {

	connStr := config.GetDatabaseURL()

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &PostgresDB{
		db: db,
	}, nil
}

func (s *PostgresDB) Init() (error, error) {
	return s.createAccountTable(), s.createBookTable()
}
