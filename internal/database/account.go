package database

import "database/sql"

func (s *PostgresDB) createAccountTable() error {
	_, err := s.db.Exec(createAccTabQ)

	return err
}

func (s *PostgresDB) CreateAccount(acc *Account) error {
	_, err := s.db.Query(
		createAccQ,
		&acc.FirstName,
		&acc.LastName,
		&acc.Email,
		&acc.CreatedAt,
		&acc.UpdatedAt,
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

func (s *PostgresDB) UpdateaAccount(*Account) error {
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
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
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
	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Email,
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
