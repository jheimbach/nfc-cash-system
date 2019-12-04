package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/JHeimbach/nfc-cash-system/server/models"
	isPkg "github.com/matryer/is"

	"net/http"
	"net/http/httptest"
	"testing"
)

type accountModelMock struct {
	create func(name, description string, startSaldo float64, groupId int, nfcChipId string) (int, error)
	list   func() ([]models.Account, error)
	read   func(id int) (models.Account, error)
	delete func(id int) error
	update func(m models.Account) error

	storage map[int]models.Account
}

func (a accountModelMock) Create(name, description string, startSaldo float64, groupId int, nfcChipId string) (int, error) {
	return a.create(name, description, startSaldo, groupId, nfcChipId)
}

func (a accountModelMock) GetAll() ([]models.Account, error) {
	return a.list()
}

func (a accountModelMock) Read(id int) (models.Account, error) {
	return a.read(id)
}

func (a accountModelMock) Delete(id int) error {
	return a.delete(id)
}

func (a accountModelMock) Update(m models.Account) error {
	return a.update(m)
}

func TestAccountService_ListAccounts(t *testing.T) {
	is := isPkg.New(t)

	tests := []struct {
		name       string
		want       []models.Account
		getAllFunc func() ([]models.Account, error)
		wantErr    bool
	}{
		{
			name: "get list with one account",
			want: genListModels(1),
			getAllFunc: func() ([]models.Account, error) {
				return genListModels(1), nil
			},
			wantErr: false,
		},
		{
			name: "get list with two accounts",
			want: genListModels(2),
			getAllFunc: func() ([]models.Account, error) {
				return genListModels(2), nil
			},
			wantErr: false,
		},
		{
			name: "get list returns an error",
			want: []models.Account{},
			getAllFunc: func() ([]models.Account, error) {
				return []models.Account{}, errors.New("error with getall")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/accounts", nil)
			accountService := AccountService{
				storage: accountModelMock{
					create: nil,
					list:   tt.getAllFunc,
				},
			}

			accountService.ListAccounts(w, r)

			if tt.wantErr {
				checkError(t, w, http.StatusInternalServerError, ErrGetAll)
				return
			}

			var accountlist []models.Account
			err := json.NewDecoder(w.Body).Decode(&accountlist)
			is.NoErr(err) // expected no error

			is.Equal(accountlist, tt.want)
		})
	}
}

func checkError(t *testing.T, w *httptest.ResponseRecorder, code int, err error) {
	t.Helper()

	if w.Code != code {
		t.Errorf("got %d status code wanted %d ", w.Code, code)
	}

	wantMsg := err.Error() + "\n"

	if w.Body.String() != wantMsg {
		t.Errorf("got err msg %q expected %q", w.Body.String(), wantMsg)
	}
}

func genListModels(num int) []models.Account {
	accounts := make([]models.Account, 0, num)

	for i := 1; i <= num; i++ {
		accounts = append(accounts, models.Account{
			ID:          i,
			Name:        "test",
			Description: "test",
			Saldo:       0,
			NfcChipId:   fmt.Sprintf("ncf_chip_%d", i),
		})
	}
	return accounts
}

func genMapModels(num int) map[int]models.Account {
	accounts := genListModels(num)
	m := make(map[int]models.Account)
	for _, acc := range accounts {
		m[acc.ID] = acc
	}

	return m
}

func TestAccountService_CreateAccount(t *testing.T) {
	is := isPkg.New(t)

	tests := []struct {
		name       string
		input      models.Account
		wantId     int
		wantErr    bool
		createFunc func(name, description string, startSaldo float64, groupId int, nfcChipId string) (int, error)
	}{
		{
			name: "create new account",
			input: models.Account{
				Name:        "tim",
				Description: "with description",
				Saldo:       120,
				NfcChipId:   "nfc_chip-id",
				GroupId:     1,
			},
			wantId:  1,
			wantErr: false,
			createFunc: func(name, description string, startSaldo float64, groupId int, nfcChipId string) (int, error) {
				return 1, nil
			},
		},
		{
			name: "create returns ErrGroupNotFound",
			input: models.Account{
				Name:        "tim",
				Description: "with description",
				Saldo:       120,
				NfcChipId:   "nfc_chip-id",
				GroupId:     10,
			},
			wantErr: true,
			createFunc: func(name, description string, startSaldo float64, groupId int, nfcChipId string) (int, error) {
				return 0, models.ErrGroupNotFound
			},
		},
		{
			name: "create returns ErrDuplicateNfcChipId",
			input: models.Account{
				Name:        "tim",
				Description: "with description",
				Saldo:       120,
				NfcChipId:   "nfc_chip-id",
				GroupId:     10,
			},
			wantErr: true,
			createFunc: func(name, description string, startSaldo float64, groupId int, nfcChipId string) (int, error) {
				return 0, models.ErrDuplicateNfcChipId
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(tt.input)
			is.NoErr(err) // marshalling error

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/accounts", &body)
			service := AccountService{
				storage: accountModelMock{
					create: tt.createFunc,
				},
			}
			service.CreateAccount(w, r)

			if tt.wantErr {
				checkError(t, w, http.StatusInternalServerError, ErrCouldNotCreateAccount)
				return
			}

			want := tt.input
			want.ID = tt.wantId
			wantBytes, err := json.Marshal(want)
			wantStr := string(wantBytes) + "\n"
			is.NoErr(err)

			if w.Body.String() != wantStr {
				t.Errorf("got %v expected %v", w.Body.String(), wantStr)
			}
		})
	}

}

func TestAccountService_GetAccount(t *testing.T) {
	is := isPkg.New(t)

	accounts := genMapModels(3)

	service := AccountService{
		storage: accountModelMock{
			read: func(id int) (models.Account, error) {
				account, ok := accounts[id]
				if !ok {
					return models.Account{}, models.ErrNotFound
				}
				return account, nil
			},
		},
	}

	tests := []struct {
		name    string
		input   int
		want    models.Account
		wantErr bool
	}{
		{
			name:    "get single account",
			input:   1,
			want:    accounts[1],
			wantErr: false,
		},
		{
			name:    "get different account",
			input:   2,
			want:    accounts[2],
			wantErr: false,
		},
		{
			name:    "get 404 account not found",
			input:   4,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/account/%d", tt.input), nil)
			service.GetAccount(w, r)

			if tt.wantErr {
				checkError(t, w, http.StatusNotFound, ErrNotFound)
				return
			}

			var got models.Account
			err := json.Unmarshal(w.Body.Bytes(), &got)
			is.NoErr(err)          // unmarshalling error
			is.Equal(got, tt.want) // got different account than expected
		})
	}

	t.Run("returns error if id is not parsable", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/account/not_a_number", nil)
		service.GetAccount(w, r)

		checkError(t, w, http.StatusInternalServerError, ErrCouldNotParseId)
		return
	})
}

func TestAccountService_DeleteAccount(t *testing.T) {

	accounts := genMapModels(3)

	service := AccountService{
		storage: accountModelMock{
			delete: func(id int) error {
				_, ok := accounts[id]
				if !ok {
					return models.ErrNotFound
				}
				delete(accounts, id)
				return nil
			},
		},
	}

	tests := []struct {
		name    string
		input   int
		want    models.Account
		wantErr bool
	}{
		{
			name:    "delete account",
			input:   1,
			wantErr: false,
		},
		{
			name:    "delete differnt account",
			input:   2,
			wantErr: false,
		},
		{
			name:    "deleting same account return 404",
			input:   2,
			wantErr: true,
		},
		{
			name:    "deleting non existent account returns 404",
			input:   4,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/account/%d", tt.input), nil)

			service.DeleteAccount(w, r)

			if tt.wantErr {
				checkError(t, w, http.StatusNotFound, ErrNotFound)
				return
			}

			// check that output is empty
			// todo maybe a success: true?
			if w.Body.String() != "" {
				t.Errorf("expected empty body, got %q", w.Body.String())
			}

			// check if account is not longer in map
			if _, ok := accounts[tt.input]; ok {
				t.Error("delete has not deleted model")
			}
		})
	}

	t.Run("returns error if id is not parsable", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, "/account/not_a_number", nil)
		service.DeleteAccount(w, r)

		checkError(t, w, http.StatusInternalServerError, ErrCouldNotParseId)
		return
	})
}
