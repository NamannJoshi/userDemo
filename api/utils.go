package api

import (
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password %v", err)
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err 
}

func StatusOK(w http.ResponseWriter, v any) error {
	return WriteJSON(w, http.StatusOK, v)
}

func BadRequest(w http.ResponseWriter, err error) error {
	return WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
}

func Unauthorized(w http.ResponseWriter, err error) error {
	return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: err.Error()})
}

func AccessForbidden(w http.ResponseWriter, err error) error {
	return WriteJSON(w, http.StatusForbidden, ApiError{Error: err.Error()})
}

func MethodNotFound(w http.ResponseWriter, err error) error {
	return WriteJSON(w, http.StatusMethodNotAllowed, ApiError{Error: err.Error()})
}

func InternalServerError(w http.ResponseWriter, err error) error {
	return WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
}