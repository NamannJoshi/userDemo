package api

import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type Account struct {
	ID              int      `json:"id"`
	FirstName       string   `json:"firstName"`
	LastName        string   `json:"lastName"`
	Email           string   `json:"email"`
	Number          int   `json:"number"`
	PasswordHash    string   `json:"passwordHash"`
}

func NewAccount(firstName, lastName, email string) *Account {
	return &Account{
		FirstName: firstName,
		LastName: lastName,
		Email: email,
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