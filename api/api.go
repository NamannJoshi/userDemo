package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)


func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error 
type ApiError struct {
	Error string
}

func makeHTTPhandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type ApiServer struct {
	listenAddr string
	store Storage
}
func NewServerApi(listenAddr string, store Storage) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		store: store,
	}
}

func (s *ApiServer)Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPhandler(s.handleDirectMethods))
	router.HandleFunc("/account/{id}", makeHTTPhandler(s.handleIDMethods))
	http.ListenAndServe(s.listenAddr, router)
}

func (s *ApiServer)handleDirectMethods(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAllAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return WriteJSON(w, http.StatusBadGateway, ApiError{Error: "method not available"})
}

func (s *ApiServer)handleIDMethods(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccountByID(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return WriteJSON(w, http.StatusBadGateway, ApiError{Error: "method not available"})
}

func (s *ApiServer)handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer)handleGetAllAccounts(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("Giggity", "Ka choda", "abs@gmail.com")
	account1 := NewAccount("Giggity", "Ka choda", "abs@gmail.com")

	accounts := []*Account{account, account1}

	type AccountResponse struct {
		ID int `json:"id"`
		Data []*Account `json:"data"`
	}

	return WriteJSON(w, http.StatusOK, AccountResponse{ID: 69, Data: accounts})
}

func (s *ApiServer)handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer)handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer)handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}