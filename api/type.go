package api

import (
	"time"
)

type CreateAccountReq struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

type UserSendRes struct {
	ID int `json:"id"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	IsAdmin bool `json:"isAdmin"`
}

type UserUpdateReq struct {
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	IsAdmin *bool `json:"isAdmin"`
	PasswordHash string `json:"passwordHash"`
}

type Account struct {
	ID              int      `json:"id"`
	FirstName       string   `json:"firstName"`
	LastName        string   `json:"lastName"`
	Email           string   `json:"email"`
	Number          int   `json:"number"`
	PasswordHash    string   `json:"passwordHash"`
	IsAdmin					bool 			`json:"isAdmin"`
	CreatedAt				time.Time `json:"createdAt"`
	UpdatedAt				time.Time `json:"updatedAt"`
}

type LoginUserReq struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRes struct {
	AccessToken string `json:"access_token"`
	User UserSendRes `json:"user"`
}

func NewAccount(firstName, lastName, email, password string) *Account {
	return &Account{
		FirstName: firstName,
		LastName: lastName,
		Email: email,
		PasswordHash: password,
		CreatedAt: time.Now().UTC(),
	}
}

type ShippingAddress struct {
	PlotNum    string `json:"plotNum"`
	Area       string `json:"area"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
}