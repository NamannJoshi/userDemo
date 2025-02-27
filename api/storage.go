package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	CreateAccountDB(*Account) (error) 
	GetAllAccountsDB() ([]*Account, error)
	GetAccountByIDDB(int) (*Account, error)
	DeleteAccountDB(int) (error)
	UpdateAccountDB(*UserUpdateReq, int) error
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
	query := `
		insert into users (firstName, lastName, email, number, password) values (
			@firstName, @lastName, @email, @number, @password
		)
	`
	args := pgx.NamedArgs{
		"firstName": account.FirstName,
		"lastName": account.LastName,
		"email": account.Email,
		"number": account.Number,
		"password": account.PasswordHash,
	}
	_, err := s.conn.Exec(context.Background(), query, args)
	return err
}

func (s *PostgreStore)GetAllAccountsDB() ([]*Account, error) {
	query := `select * from users;`

	rows, err := s.conn.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("error occured while fetching accounts from DB: %v", err)
	}
	defer rows.Close()

	var users []*Account
	for rows.Next() {
			var user Account
			err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Number, &user.PasswordHash)

			if err != nil {
			return nil, fmt.Errorf("error occured while Scanning accounts from DB: %v", err)
		}
		users = append(users, &user)
	}
	return users, nil
	
}

func(s *PostgreStore) GetAccountByIDDB(accountId int) (*Account, error) {
	query := `
		select * from users where id = @accountId;
	`
	args := pgx.NamedArgs{
		"accountId": strconv.Itoa(accountId),
	}
	row := s.conn.QueryRow(context.Background(), query, args)

	var res Account 
	err := row.Scan(&res.ID, &res.FirstName, &res.LastName, &res.Email, &res.Number, &res.PasswordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("row not found")
		}
		return nil, err
	}
	
	return &res, nil
}

func (s *PostgreStore)DeleteAccountDB(accountID int) (error) {
	query := `
		delete from users where id = @accountID;
	`
	args := pgx.NamedArgs{
		"accountId": strconv.Itoa(accountID),
	}
	_, err := s.conn.Exec(context.Background(), query, args) 
	if err != nil {
		return fmt.Errorf("account %d not found", accountID)
	}
	return err
}

func (s *PostgreStore)UpdateAccountDB(account *UserUpdateReq, accountID int) (error) {

	//change updates data type of map's value and its if statement and data type in "for" loop's sprintf
	updates := make(map[string]string)
	if account.FirstName != "" {
		updates["firstName"] = account.FirstName
	}
	if account.LastName != "" {
		updates["lastName"] = account.LastName
	}
	if account.Email != "" {
		updates["email"] = account.Email
	}
	if account.PasswordHash != "" {
		updates["passwordHash"] = account.PasswordHash
	}

	if len(updates) == 0 {
		return nil
	}

	var setClause []string
	args := pgx.NamedArgs{}

	for key, value :=range updates {
		setClause = append(setClause, fmt.Sprintf("%s = @%s", key, key))
		args[key] = value
	}
	args["accountID"] = accountID
	query := fmt.Sprintf("update users set %s where id = @accountID", strings.Join(setClause, ", "))

	fmt.Println(query)
	fmt.Println(args)
	
	_, err := s.conn.Exec(context.Background(), query, args)
	if err != nil {
		return fmt.Errorf("error updating the users profile")
	}
	return nil
}

