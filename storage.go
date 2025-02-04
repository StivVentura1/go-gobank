package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CraeteAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
	//GetAccountByNumber(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=******* host=127.0.0.1 port=5433 sslmode=disable options='-c search_path=go-gobank'"
	db, err := sql.Open("postgres", connStr)
	log.Println("JSON API server running on port: ", db.Stats())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}
func (s *PostgresStore) init() error {
	return s.createAccountTable()
}
func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    number BIGINT UNIQUE,
    encrypted_password VARCHAR(100),
    balance NUMERIC(10, 2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`

	_, err := s.db.Exec(query)
	return err
}
func (s *PostgresStore) CraeteAccount(acc *Account) error {
	query := `insert into accounts (first_name,last_name,number,balance,created_at) values ($1,$2,$3,$4,$5)`
	_, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
func (ps *PostgresStore) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	_, err := s.db.Query("delete from accounts where id = $1", id)
	return err

}
func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select *from accounts")
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		/*err := rows.Scan(
				&account.ID,
				&account.FirstName,
				&account.LastName,
				&account.Number,
				&account.Balance,
				&account.CreatedAt)

			if err != nil {
				return nil, err
			}
			accounts = append(accounts, account)
		}*/
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, err
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query("select *from accounts where id = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		//&account.EncryptedPassword,
		&account.Balance,
		&account.CreatedAt)

	return account, err
}
