package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/JHeimbach/nfc-cash-system/server/models"
)

var (
	ErrGetAll                = errors.New("could not load list of accounts")
	ErrCouldNotCreateAccount = errors.New("could not save new account")
	ErrNotFound              = errors.New("could not find account")
	ErrCouldNotParseId       = errors.New("could not parse account id")
	ErrSomethingWentWrong    = errors.New("something went wrong")
)

type AccountService struct {
	storage models.AccountStorager
}

func (a AccountService) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := a.storage.GetAll()
	if err != nil {
		http.Error(w, ErrGetAll.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(accounts)
}

func (a AccountService) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	json.NewDecoder(r.Body).Decode(&account)

	id, err := a.storage.Create(account.Name, account.Description, account.Saldo, account.ID, account.NfcChipId)
	if err != nil {
		http.Error(w, ErrCouldNotCreateAccount.Error(), http.StatusInternalServerError)
		return
	}

	account.ID = id
	json.NewEncoder(w).Encode(account)
}

func (a AccountService) GetAccount(w http.ResponseWriter, r *http.Request) {
	id, err := getAccountId(r.URL.Path)
	if err != nil {
		http.Error(w, ErrCouldNotParseId.Error(), http.StatusInternalServerError)
		return
	}

	account, err := a.storage.Read(id)

	if err != nil {
		http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(account)
}

func (a AccountService) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id, err := getAccountId(r.URL.Path)
	if err != nil {
		http.Error(w, ErrCouldNotParseId.Error(), http.StatusInternalServerError)
		return
	}
	err = a.storage.Delete(id)

	if err != nil {
		if err == models.ErrNotFound {
			http.Error(w, ErrNotFound.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, ErrSomethingWentWrong.Error(), http.StatusNotFound)
		return
	}
}

func (a AccountService) UpdateAccount(w http.ResponseWriter, r *http.Request) {

}

func getAccountId(path string) (int, error) {
	return strconv.Atoi(strings.Replace(path, "/account/", "", -1))
}
