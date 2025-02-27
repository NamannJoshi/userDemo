package api

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	CreateAccountDB(*Account) error 
	GetAllAccountsDB() ([]*Account, error)
	GetAccountByIDDB(int) (*Account, error)
	DeleteAccountDB(int) error
	UpdateAccountDB(*Account, int) error
}

type PostgreStore struct {
	conn *pgx.Conn
}

func NewPostgreStore() (*PostgreStore, error) {
	connString := "postgres://postgres:mukeshakun@localhost:5432/test"

	conn, err := pgx.Connect(context.Background(), connString)

	return &PostgreStore{
		conn: conn,
	}, err
}

func (s *PostgreStore)Init() error {
	return s.CreateAccountTable()
}

func (s *PostgreStore)CreateAccountTable() error {
	query := `
		create table if not exists users (
			id serial primary key,
			firstName varchar(100) not null,
			lastName varchar(100) not null,
			email varchar(200) not null,
			number int not null,
			password varchar(200) not null
		)
	`
	_, err := s.conn.Exec(context.Background(), query)
	return err
}

func (s *PostgreStore)CreateAccountDB(account *Account) (error) {
	return nil
}

func (s *PostgreStore)GetAllAccountsDB() ([]*Account, error) {
	return nil, nil
}

func (s *PostgreStore)GetAccountByIDDB(accountID int) (*Account, error) {
	return nil, nil
}

func (s *PostgreStore)DeleteAccountDB(accountID int) (error) {
	return nil
}

func (s *PostgreStore)UpdateAccountDB(account *Account, accountID int) (error) {
	return nil
}

