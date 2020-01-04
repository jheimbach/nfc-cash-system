package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test"
	isPkg "github.com/matryer/is"
)

func TestAccountserver_E2E_ListAccounts(t *testing.T) {
	test.IsIntegrationTest(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode   int
		errMsg       string
		accountsLen  int
		accountTotal int32
	}
	tests := []struct {
		name         string
		accessToken  string
		want         want
		pagingLimit  int
		pagingOffset int
		groupId      int
	}{
		{
			name: "no accesstoken given",
			want: want{
				statusCode: http.StatusUnauthorized,
				errMsg:     "authorization header required",
			},
		},
		{
			name:        "get all accounts",
			accessToken: aTkn,
			want: want{
				statusCode:   http.StatusOK,
				accountsLen:  100,
				accountTotal: 100,
			},
		},
		{
			name:        "get first 10 accounts",
			accessToken: aTkn,
			pagingLimit: 10,
			want: want{
				statusCode:   http.StatusOK,
				accountsLen:  10,
				accountTotal: 100,
			},
		},
		{
			name:         "get second 10 accounts",
			accessToken:  aTkn,
			pagingLimit:  10,
			pagingOffset: 10,
			want: want{
				statusCode:   http.StatusOK,
				accountsLen:  10,
				accountTotal: 100,
			},
		},
		{
			name:        "filter by group",
			accessToken: aTkn,
			groupId:     1,
			want: want{
				statusCode:   http.StatusOK,
				accountsLen:  10,
				accountTotal: 10,
			},
		},
		{
			name:        "filter by group with limit",
			accessToken: aTkn,
			groupId:     1,
			pagingLimit: 5,
			want: want{
				statusCode:   http.StatusOK,
				accountsLen:  5,
				accountTotal: 10,
			},
		},
		{
			name:         "filter by group with limit and offset",
			accessToken:  aTkn,
			groupId:      1,
			pagingLimit:  5,
			pagingOffset: 8,
			want: want{
				statusCode:   http.StatusOK,
				accountsLen:  2,
				accountTotal: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := url.Parse(RestUrlWithPath("v1/accounts"))
			if err != nil {
				t.Fatalf("could not parse url: %s; %v", RestUrlWithPath("v1/accounts"), err)
			}

			if tt.pagingLimit != 0 || tt.groupId != 0 {
				q := path.Query()
				if tt.groupId != 0 {
					q.Add("group_id", strconv.Itoa(tt.groupId))
				}
				if tt.pagingLimit != 0 {
					q.Add("paging.limit", strconv.Itoa(tt.pagingLimit))
				}
				if tt.pagingOffset != 0 {
					q.Add("paging.offset", strconv.Itoa(tt.pagingOffset))
				}
				path.RawQuery = q.Encode()
			}

			req, err := http.NewRequest(http.MethodGet, path.String(), nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			if tt.accessToken != "" {
				req.Header.Add("Authorization", "Bearer "+tt.accessToken)
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("could not request accounts: %v", err)
			}
			defer res.Body.Close()

			if tt.want.statusCode != http.StatusOK {
				checkError(t, res, tt.want.statusCode, tt.want.errMsg)
				return
			}

			var accounts api.ListAccountsResponse
			err = json.NewDecoder(res.Body).Decode(&accounts)

			if err != nil {
				t.Fatalf("could not parse accounts: %v", err)
			}

			if l := len(accounts.Accounts); l != tt.want.accountsLen {
				t.Errorf("got %d accounts, wanted %d", l, tt.want.accountsLen)
			}

			if accounts.TotalCount != tt.want.accountTotal {
				t.Errorf("got totalcount %d, wanted %d", accounts.TotalCount, tt.want.accountTotal)
			}
		})
	}
}

func TestAccountserver_E2E_GetAccount(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode int
		errMsg     string
		account    api.Account
	}
	tests := []struct {
		name        string
		accessToken string
		accountId   int
		want        want
	}{
		{
			name:      "no accesstoken given",
			accountId: 1,
			want: want{
				statusCode: http.StatusUnauthorized,
				errMsg:     "authorization header required",
			},
		},
		{
			name:        "get account with id 1",
			accessToken: aTkn,
			accountId:   1,
			want: want{
				statusCode: http.StatusOK,
				account: api.Account{
					Id:          1,
					Name:        "Laverne Blackstock",
					Description: "Itchy Eye",
					Saldo:       436,
					NfcChipId:   "Hv8mnajqzIKO",
					Group: &api.Group{
						Id:          7,
						Name:        "PSS World Medical, Inc.",
						Description: "",
						CanOverdraw: true,
					},
				},
			},
		},
		{
			name:        "account with invalid id",
			accessToken: aTkn,
			accountId:   -45,
			want: want{
				statusCode: http.StatusNotFound,
				errMsg:     "could not find account",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			req, err := http.NewRequest(http.MethodGet, RestUrlWithPath(fmt.Sprintf("v1/account/%d", tt.accountId)), nil)
			is.NoErr(err) // could not create request
			if tt.accessToken != "" {
				req.Header.Add("Authorization", "Bearer "+tt.accessToken)
			}

			res, err := http.DefaultClient.Do(req)
			is.NoErr(err) // request failed
			defer res.Body.Close()

			if tt.want.statusCode != http.StatusOK {
				checkError(t, res, tt.want.statusCode, tt.want.errMsg)
				return
			}

			var account api.Account
			err = json.NewDecoder(res.Body).Decode(&account)
			is.NoErr(err) // could not decode account

			is.Equal(account, tt.want.account) // account is not the expected
		})
	}
}

func TestAccountserver_E2E_CreateAccount(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode int
		errMsg     string
		account    api.Account
	}
	tests := []struct {
		name        string
		accessToken string
		body        *api.CreateAccountRequest
		want        want
	}{
		{
			name: "no accesstoken given",
			want: want{
				statusCode: http.StatusUnauthorized,
				errMsg:     "authorization header required",
			},
		},
		{
			name:        "create new account",
			accessToken: aTkn,
			body: &api.CreateAccountRequest{
				Name:        "test account",
				Description: "for testing",
				Saldo:       1000,
				NfcChipId:   "t3stch1p",
				GroupId:     1,
			},
			want: want{
				statusCode: http.StatusOK,
				account: api.Account{
					Id:          101,
					Name:        "test account",
					Description: "for testing",
					Saldo:       1000,
					NfcChipId:   "t3stch1p",
					Group: &api.Group{
						Id:   1,
						Name: "H2O Plus",
					},
				},
			},
		},
		{
			name:        "create new account with used nfc chip id",
			accessToken: aTkn,
			body: &api.CreateAccountRequest{
				Name:        "test account",
				Description: "for testing",
				Saldo:       1000,
				NfcChipId:   "0XPPQy4ZkO7",
				GroupId:     1,
			},
			want: want{
				statusCode: http.StatusConflict,
				errMsg:     "nfc chip is already in use",
			},
		},
		{
			name:        "create new account with nonexistent group",
			accessToken: aTkn,
			body: &api.CreateAccountRequest{
				Name:        "test account",
				Description: "for testing",
				Saldo:       1000,
				NfcChipId:   "ofzGN0eS34",
				GroupId:     -45,
			},
			want: want{
				statusCode: http.StatusNotFound,
				errMsg:     "could not find group",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			body, err := json.Marshal(tt.body)
			is.NoErr(err) // could not marshal body

			req, err := http.NewRequest(http.MethodPost, RestUrlWithPath("v1/accounts"), bytes.NewReader(body))
			is.NoErr(err) // could not create request

			if tt.accessToken != "" {
				req.Header.Add("Authorization", "Bearer "+tt.accessToken)
			}

			res, err := http.DefaultClient.Do(req)
			is.NoErr(err) // request failed
			defer res.Body.Close()

			if tt.want.statusCode != http.StatusOK {
				checkError(t, res, tt.want.statusCode, tt.want.errMsg)
				return
			}

			var account api.Account
			err = json.NewDecoder(res.Body).Decode(&account)
			is.NoErr(err) // could not decode account

			is.Equal(account, tt.want.account) // account is not the expected
		})
	}
}

func TestAccountserver_E2E_UpdateAccount(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode int
		errMsg     string
		account    api.Account
	}
	tests := []struct {
		name        string
		accessToken string
		body        *api.Account
		want        want
	}{
		{
			name: "no accesstoken given",
			body: &api.Account{
				Id: 1,
			},
			want: want{
				statusCode: http.StatusUnauthorized,
				errMsg:     "authorization header required",
			},
		},
		{
			name:        "update account name",
			accessToken: aTkn,
			body: &api.Account{
				Id:          1,
				Name:        "Laverne",
				Description: "Itchy Eye",
				Saldo:       436,
				NfcChipId:   "Hv8mnajqzIKO",
				Group: &api.Group{
					Id: 7,
				},
			},
			want: want{
				statusCode: http.StatusOK,
				account: api.Account{
					Id:          1,
					Name:        "Laverne",
					Description: "Itchy Eye",
					Saldo:       436,
					NfcChipId:   "Hv8mnajqzIKO",
					Group: &api.Group{
						Id:          7,
						Name:        "PSS World Medical, Inc.",
						CanOverdraw: true,
					},
				},
			},
		},
		{
			name:        "update account description",
			accessToken: aTkn,
			body: &api.Account{
				Id:          1,
				Name:        "Laverne",
				Description: "",
				Saldo:       436,
				NfcChipId:   "Hv8mnajqzIKO",
				Group: &api.Group{
					Id: 7,
				},
			},
			want: want{
				statusCode: http.StatusOK,
				account: api.Account{
					Id:          1,
					Name:        "Laverne",
					Description: "",
					Saldo:       436,
					NfcChipId:   "Hv8mnajqzIKO",
					Group: &api.Group{
						Id:          7,
						Name:        "PSS World Medical, Inc.",
						CanOverdraw: true,
					},
				},
			},
		},
		{
			name:        "try to update account saldo",
			accessToken: aTkn,
			body: &api.Account{
				Id:          1,
				Name:        "Laverne",
				Description: "",
				Saldo:       1000,
				NfcChipId:   "Hv8mnajqzIKO",
				Group: &api.Group{
					Id: 7,
				},
			},
			want: want{
				statusCode: http.StatusForbidden,
				errMsg:     "can not update account saldo trough update",
			},
		},
		{
			name:        "move acccount to different group",
			accessToken: aTkn,
			body: &api.Account{
				Id:          1,
				Name:        "Laverne",
				Description: "",
				Saldo:       436,
				NfcChipId:   "Hv8mnajqzIKO",
				Group: &api.Group{
					Id: 1,
				},
			},
			want: want{
				statusCode: http.StatusOK,
				account: api.Account{
					Id:          1,
					Name:        "Laverne",
					Description: "",
					Saldo:       436,
					NfcChipId:   "Hv8mnajqzIKO",
					Group: &api.Group{
						Id:   1,
						Name: "H2O Plus",
					},
				},
			},
		},
		{
			name:        "try to update to non existant group",
			accessToken: aTkn,
			body: &api.Account{
				Id:          1,
				Name:        "Laverne",
				Description: "",
				Saldo:       436,
				NfcChipId:   "Hv8mnajqzIKO",
				Group: &api.Group{
					Id: -45,
				},
			},
			want: want{
				statusCode: http.StatusNotFound,
				errMsg:     "could not find group",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			body, err := json.Marshal(tt.body)
			is.NoErr(err) // could not marshal body

			req, err := http.NewRequest(http.MethodPut, RestUrlWithPath(fmt.Sprintf("v1/account/%d", tt.body.Id)), bytes.NewReader(body))
			is.NoErr(err) // could not create request

			if tt.accessToken != "" {
				req.Header.Add("Authorization", "Bearer "+tt.accessToken)
			}

			res, err := http.DefaultClient.Do(req)
			is.NoErr(err) // request failed
			defer res.Body.Close()

			if tt.want.statusCode != http.StatusOK {
				checkError(t, res, tt.want.statusCode, tt.want.errMsg)
				return
			}

			var account api.Account
			err = json.NewDecoder(res.Body).Decode(&account)
			is.NoErr(err) // could not decode account

			is.Equal(account, tt.want.account) // account is not the expected
		})
	}
}
func TestAccountserver_E2E_DeleteAccount(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	tests := []struct {
		name        string
		accessToken string
		accountId   int
		statusCode  int
		errMsg      string
	}{
		{
			name:       "no accesstoken given",
			accountId:  1,
			statusCode: http.StatusUnauthorized,
			errMsg:     "authorization header required",
		},
		{
			name:        "delete account with id 1",
			accessToken: aTkn,
			accountId:   1,
			statusCode:  http.StatusOK,
		},
		{
			name:        "delete account with invalid id",
			accessToken: aTkn,
			accountId:   -45,
			statusCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			req, err := http.NewRequest(http.MethodDelete, RestUrlWithPath(fmt.Sprintf("v1/account/%d", tt.accountId)), nil)
			is.NoErr(err) // could not create request

			if tt.accessToken != "" {
				req.Header.Add("Authorization", "Bearer "+tt.accessToken)
			}

			res, err := http.DefaultClient.Do(req)
			is.NoErr(err) // request failed
			defer res.Body.Close()

			if tt.statusCode != http.StatusOK {
				checkError(t, res, tt.statusCode, tt.errMsg)
				return
			}
			b, err := ioutil.ReadAll(res.Body)
			is.NoErr(err) // could not read body

			if string(b) != "{}" {
				t.Errorf("expected empty body, got %q", b)
			}
		})
	}
}

func checkError(t *testing.T, response *http.Response, code int, errMsg string) {
	t.Helper()

	if response.StatusCode != code {
		t.Errorf("got statuscode %d, expected %d", response.StatusCode, code)
	}

	var jsonErr *api.Status
	err := json.NewDecoder(response.Body).Decode(&jsonErr)
	if err != nil {
		t.Fatalf("could not parse error: %v", err)
	}

	if jsonErr.Message != errMsg {
		t.Errorf("got err msg: %q, wanted: %q", jsonErr.Message, errMsg)
	}
}
