package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func (s *PostgresDB) createAccountTable() error {
	_, err := s.db.Exec(createAccTabQ)

	return err
}

func (s *PostgresDB) CreateAccount(account *Account) error {
	_, err := s.db.Query(
		createAccQ,
		&account.FirstName,
		&account.LastName,
		&account.Email,
		&account.Password,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) DeleteAccount(id string) error {
	_, err := s.db.Exec(deleteAccQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateAccount(id string, account *Account) error {
	accId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}

	_, err := s.db.Query(
		updateAccQ,
		accId,
		&account.FirstName,
		&account.LastName,
		&account.Email,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetAccountById(id string) (*Account, error) {
	row := s.db.QueryRow(getAccQ, id)

	var account Account

	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Email,
		&account.Password,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Account with uuid '%v' does not exist", id)
		}
		return nil, err
	}

	return &account, nil
}

func (s *PostgresDB) GetAccounts() (*[]*Account, error) {
	rows, err := s.db.Query(getAccsQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Email,
			&account.Password,
			&account.CreatedAt,
			&account.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}
	return &accounts, nil
}
