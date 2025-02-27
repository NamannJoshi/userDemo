package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func getID(r *http.Request) (int, error) {
	vars := mux.Vars(r)["id"]
	id, err := strconv.Atoi(vars)
	if err != nil {
		return -1, err
	}
	return id, nil
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
	fmt.Println("get by id endpoint working!")
	id, err := getID(r)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "invalid request"})
	}
	
	account, err := s.store.GetAccountByIDDB(id)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (s *ApiServer)handleGetAllAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAllAccountsDB()
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error:"invalid request"})
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *ApiServer)handleCreateAccount(w http.ResponseWriter, r *http.Request) (error) {
	fmt.Println("Create api endpoint is working")
	var req CreateAccountReq
	json.NewDecoder(r.Body).Decode(&req)

	account := NewAccount(req.FirstName, req.LastName, req.Email)
	if err := s.store.CreateAccountDB(account); err != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "invalid request"}) 
	}

	return WriteJSON(w, http.StatusOK, &UserSendRes{account.ID, account.FirstName, account.LastName})
}

func (s *ApiServer)handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
	}

	errar := s.store.DeleteAccountDB(id)
	if errar != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: errar.Error()})
	}
	return WriteJSON(w, http.StatusOK, map[string]string{"message":"account deleted successfully"})
}

func (s *ApiServer)handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return fmt.Errorf("error while fetching id: %d", id)
	}
	var req *UserUpdateReq
	json.NewDecoder(r.Body).Decode(&req)

	fmt.Println(req)

	errar := s.store.UpdateAccountDB(req, id); 
	if errar != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: errar.Error()})
	}

	fmt.Println(req)
	return WriteJSON(w, http.StatusOK, map[string]string{"message": "account updated successfully"})
}