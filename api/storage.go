package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Storage interface {
	CreateAccountDB(*Account) (error) 
	GetAllAccountsDB() ([]*Account, error)
	GetAccountByIDDB(int) (*Account, error)
	GetAccountByEmailDB(string) (*Account, error)
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

func (s *PostgreStore)DropTable() error {
	query := `drop table users;`
	_, err := s.conn.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("table deletion can't be handled")
	}
	return nil
}

func (s * PostgreStore) CreateSessions() error {
	query := `
		create table sessions (
			id varchar(255) primary key not null,
			email varchar(255) not null,
			refreshToken varchar(512) not null,
			isRevoked bool not null default false,
			createdAt datetime default (not()),
			expiresAt datetime
		);
	`
	_, err := s.conn.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("error creating sessions table")
	}
	return nil
}
func (s *PostgreStore)DropSessionTable() error {
	query := `
		drop table sessions;
	`
	_, err := s.conn.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("session table couldn't get deleted")
	}
	return nil
}

func (s *PostgreStore)CreateAccountTable() error {
	query := `
		create table if not exists users (
			id serial primary key,
			firstName varchar(100) not null,
			lastName varchar(100) not null,
			email varchar(200) not null unique,
			number int not null,
			passwordHash varchar(200) not null,
			isAdmin boolean not null,
			createdAt datetime not null,
			updatedAt datetime not null
		)
	`
	_, err := s.conn.Exec(context.Background(), query)
	return err
}

func (s *PostgreStore)CreateAccountDB(account *Account) (error) {
	query := `
		insert into users (firstName, lastName, email, number, passwordHash, createdAt, updatedAt, isAdmin) values (
			@firstName, @lastName, @email, @number, @passwordHash, @createdAt, @updatedAt, @isAdmin
		)
	`
	args := pgx.NamedArgs{
		"firstName": account.FirstName,
		"lastName": account.LastName,
		"email": account.Email,
		"number": account.Number,
		"passwordHash": account.PasswordHash,
		"createdAt": account.CreatedAt,
		"updatedAt": account.CreatedAt,
		"isAdmin": account.IsAdmin,
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
			err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Number, &user.PasswordHash, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

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
	err := row.Scan(&res.ID, &res.FirstName, &res.LastName, &res.Email, &res.Number, &res.PasswordHash, &res.IsAdmin, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, err
	}
	
	return &res, nil
}
func(s *PostgreStore) GetAccountByEmailDB(accountId string) (*Account, error) {
	query := `
		select * from users where email = @accountId;
	`
	args := pgx.NamedArgs{
		"accountId": (accountId),
	}
	row := s.conn.QueryRow(context.Background(), query, args)

	var res Account 
	err := row.Scan(&res.ID, &res.FirstName, &res.LastName, &res.Email, &res.Number, &res.PasswordHash, &res.IsAdmin, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("account not found")
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
	updates := make(map[string]any)
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
		hashed, err := HashPassword(account.PasswordHash)
		if err != nil {
			panic(err)
		}
		updates["passwordHash"] = hashed
	}
	if account.IsAdmin != nil {
		updates["isAdmin"] = *account.IsAdmin
	}
	updates["updatedAt"] = time.Now().UTC().Format(time.ANSIC)

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
	
	_, err := s.conn.Exec(context.Background(), query, args)
	if err != nil {
		return fmt.Errorf("error updating the users profile")
	}
	return nil
}

