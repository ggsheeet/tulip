package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *PostgresDB) createAccountTable() error {
	_, err := s.db.Exec(createAccTabQ)

	return err
}

func (s *PostgresDB) CreateAccount(account *Account) (uuid.UUID, error) {
	err := s.db.QueryRow(
		createAccQ,
		&account.FirstName,
		&account.LastName,
		&account.Email,
		&account.Phone,
		&account.CreatedAt,
		&account.UpdatedAt,
	).Scan(&account.ID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("error: %v", err)
	}

	return account.ID, nil
}

func (s *PostgresDB) DeleteAccount(id uuid.UUID) error {
	_, err := s.db.Exec(deleteAccQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateAccount(id uuid.UUID, account *Account) error {
	_, err := s.db.Query(
		updateAccQ,
		id,
		&account.FirstName,
		&account.LastName,
		&account.Email,
		&account.Phone,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetAccountById(id uuid.UUID) (*Account, error) {
	row := s.db.QueryRow(getAccQ, id)

	var account Account

	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Email,
		&account.Phone,
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

func (s *PostgresDB) GetAccountByEmail(email string) (*Account, error) {
	row := s.db.QueryRow(getAccByEmailQ, email)

	var account Account

	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Email,
		&account.Phone,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Account with email '%v' does not exist", email)
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
			&account.Phone,
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
