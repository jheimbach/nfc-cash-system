package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	isPkg "github.com/matryer/is"
)

func TestTransactionServer_E2E_ListTransactions(t *testing.T) {
	test.IsIntegrationTest(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode        int
		errMsg            string
		transactionsLen   int
		transactionsTotal int32
	}
	tests := []struct {
		name         string
		accessToken  string
		want         want
		pagingLimit  int
		pagingOffset int
		order        string
	}{
		{
			name: "no accesstoken given",
			want: want{
				statusCode: http.StatusUnauthorized,
				errMsg:     "authorization header required",
			},
		},
		{
			name:        "get all transactions",
			accessToken: aTkn,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   1000,
				transactionsTotal: 1000,
			},
		},
		{
			name:        "get first 10 transactions",
			accessToken: aTkn,
			pagingLimit: 10,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   10,
				transactionsTotal: 1000,
			},
		},
		{
			name:         "get second 10 transactions",
			accessToken:  aTkn,
			pagingLimit:  10,
			pagingOffset: 10,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   10,
				transactionsTotal: 1000,
			},
		},
		{
			name:        "order desc",
			accessToken: aTkn,
			order:       "desc",
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   1000,
				transactionsTotal: 1000,
			},
		},
		{
			name:        "order asc",
			accessToken: aTkn,
			order:       "asc",
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   1000,
				transactionsTotal: 1000,
			},
		},
		{
			name:        "order by desc with limit",
			accessToken: aTkn,
			order:       "desc",
			pagingLimit: 5,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   5,
				transactionsTotal: 1000,
			},
		},
		{
			name:        "order by asc with limit",
			accessToken: aTkn,
			order:       "asc",
			pagingLimit: 5,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   5,
				transactionsTotal: 1000,
			},
		},
		{
			name:         "order by desc with limit and offset",
			accessToken:  aTkn,
			order:        "desc",
			pagingLimit:  5,
			pagingOffset: 8,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   5,
				transactionsTotal: 1000,
			},
		},
		{
			name:         "order by asc with limit and offset",
			accessToken:  aTkn,
			order:        "asc",
			pagingLimit:  5,
			pagingOffset: 8,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   5,
				transactionsTotal: 1000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := url.Parse(RestUrlWithPath("v1/transactions"))
			if err != nil {
				t.Fatalf("could not parse url: %s; %v", RestUrlWithPath("v1/transactions"), err)
			}

			if tt.pagingLimit != 0 || tt.order != "" {
				q := path.Query()
				if tt.order != "" {
					q.Add("order", tt.order)
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
				t.Fatalf("could not request transactionsResponse: %v", err)
			}
			defer res.Body.Close()

			if tt.want.statusCode != http.StatusOK {
				checkError(t, res, tt.want.statusCode, tt.want.errMsg)
				return
			}

			var transactionsResponse api.ListTransactionsResponse
			err = jsonpb.Unmarshal(res.Body, &transactionsResponse)

			if err != nil {
				t.Fatalf("could not parse transactionsResponse: %v", err)
			}
			checkTransactionLen(t, transactionsResponse, tt.want.transactionsLen)
			checkTotalCount(t, transactionsResponse, tt.want.transactionsTotal)

			checkOrder(t, transactionsResponse, tt.order)
		})
	}
}

func TestTransactionServer_E2E_ListTransactionsByAccount(t *testing.T) {
	test.IsIntegrationTest(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode        int
		errMsg            string
		transactionsLen   int
		transactionsTotal int32
	}
	tests := []struct {
		name         string
		accessToken  string
		want         want
		pagingLimit  int
		pagingOffset int
		order        string
		accountID    int32
	}{
		{
			name: "no accesstoken given",
			want: want{
				statusCode: http.StatusUnauthorized,
				errMsg:     "authorization header required",
			},
		},
		{
			name:        "get all transactions for account id 1",
			accessToken: aTkn,
			accountID:   1,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   10,
				transactionsTotal: 10,
			},
		},
		{
			name:        "get first 10 transactions",
			accessToken: aTkn,
			accountID:   1,
			pagingLimit: 10,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   10,
				transactionsTotal: 10,
			},
		},
		{
			name:         "get second 5 transactions",
			accessToken:  aTkn,
			accountID:    1,
			pagingLimit:  10,
			pagingOffset: 5,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   5,
				transactionsTotal: 10,
			},
		},
		{
			name:        "order desc",
			accessToken: aTkn,
			accountID:   1,
			order:       "desc",
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   10,
				transactionsTotal: 10,
			},
		},
		{
			name:        "order asc",
			accessToken: aTkn,
			accountID:   1,
			order:       "asc",
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   10,
				transactionsTotal: 10,
			},
		},
		{
			name:        "order by desc with limit",
			accessToken: aTkn,
			accountID:   1,
			order:       "desc",
			pagingLimit: 5,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   5,
				transactionsTotal: 10,
			},
		},
		{
			name:        "order by asc with limit",
			accessToken: aTkn,
			accountID:   1,
			order:       "asc",
			pagingLimit: 5,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   5,
				transactionsTotal: 10,
			},
		},
		{
			name:         "order by desc with limit and offset",
			accessToken:  aTkn,
			accountID:    1,
			order:        "desc",
			pagingLimit:  5,
			pagingOffset: 8,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   2,
				transactionsTotal: 10,
			},
		},
		{
			name:         "order by asc with limit and offset",
			accessToken:  aTkn,
			accountID:    1,
			order:        "asc",
			pagingLimit:  5,
			pagingOffset: 8,
			want: want{
				statusCode:        http.StatusOK,
				transactionsLen:   2,
				transactionsTotal: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rPath := RestUrlWithPath(fmt.Sprintf("v1/account/%d/transactions", tt.accountID))
			path, err := url.Parse(rPath)
			if err != nil {
				t.Fatalf("could not parse url: %s; %v", rPath, err)
			}

			if tt.pagingLimit != 0 || tt.order != "" {
				q := path.Query()
				if tt.order != "" {
					q.Add("order", tt.order)
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
				t.Fatalf("could not request transactionsResponse: %v", err)
			}
			defer res.Body.Close()

			if tt.want.statusCode != http.StatusOK {
				checkError(t, res, tt.want.statusCode, tt.want.errMsg)
				return
			}

			var transactionsResponse api.ListTransactionsResponse
			err = jsonpb.Unmarshal(res.Body, &transactionsResponse)

			if err != nil {
				t.Fatalf("could not parse transactionsResponse: %v", err)
			}

			checkTransactionLen(t, transactionsResponse, tt.want.transactionsLen)
			checkTotalCount(t, transactionsResponse, tt.want.transactionsTotal)
			checkOrder(t, transactionsResponse, tt.order)
			checkAccountId(t, transactionsResponse, tt.accountID)
		})
	}
}

func TestTransactionServer_E2E_GetTransaction(t *testing.T) {
	test.IsIntegrationTest(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode  int
		errMsg      string
		transaction api.Transaction
	}
	tests := []struct {
		name          string
		accessToken   string
		want          want
		accountID     int32
		transactionID int32
	}{
		{
			name: "no accesstoken given",
			want: want{
				statusCode: http.StatusUnauthorized,
				errMsg:     "authorization header required",
			},
		},
		{
			name:          "get transaction with id 1 and account id 20",
			accessToken:   aTkn,
			accountID:     20,
			transactionID: 1,
			want: want{
				statusCode: http.StatusOK,
				transaction: api.Transaction{
					Id:       1,
					OldSaldo: 540,
					NewSaldo: 539,
					Amount:   1,
					Created: func() *timestamp.Timestamp {
						t, _ := ptypes.TimestampProto(time.Date(2018, 12, 10, 1, 58, 6, 0, time.UTC))
						return t
					}(),
					Account: &api.Account{
						Id:          20,
						Name:        "Florida Duesberry",
						Description: "",
						Saldo:       495,
						NfcChipId:   "rKlqNQsxt",
						Group: &api.Group{
							Id:   9,
							Name: "Pharmacia and Upjohn Company",
						},
					},
				},
			},
		},
		{
			name:          "get transaction with id 1 and account id 2 returns 404 err",
			accessToken:   aTkn,
			accountID:     2,
			transactionID: 1,
			want: want{
				statusCode: http.StatusNotFound,
				errMsg:     "could not find transaction",
			},
		},
		{
			name:          "transaction with id 1111 and account id 2 returns 404 err",
			accessToken:   aTkn,
			accountID:     2,
			transactionID: 1111,
			want: want{
				statusCode: http.StatusNotFound,
				errMsg:     "could not find transaction",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rPath := RestUrlWithPath(fmt.Sprintf("v1/account/%d/transactions/%d", tt.accountID, tt.transactionID))
			path, err := url.Parse(rPath)
			if err != nil {
				t.Fatalf("could not parse url: %s; %v", rPath, err)
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
				t.Fatalf("could not request transaction: %v", err)
			}
			defer res.Body.Close()

			if tt.want.statusCode != http.StatusOK {
				checkError(t, res, tt.want.statusCode, tt.want.errMsg)
				return
			}

			var transaction api.Transaction
			err = jsonpb.Unmarshal(res.Body, &transaction)

			if err != nil {
				t.Fatalf("could not parse transaction: %v", err)
			}
			if !reflect.DeepEqual(transaction, tt.want.transaction) {
				t.Errorf("got %v, wanted %v", transaction, tt.want.transaction)
			}
		})
	}
}
func TestTransactionServer_E2E_CreateTransaction(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode int
		errMsg     string
		account    api.Transaction
	}
	tests := []struct {
		name        string
		accessToken string
		body        api.CreateTransactionRequest
		want        want
	}{
		{
			name: "no accesstoken given",
			body: api.CreateTransactionRequest{
				AccountId: 1,
			},
			want: want{
				statusCode: http.StatusUnauthorized,
				errMsg:     "authorization header required",
			},
		},
		{
			name:        "create new transaction",
			accessToken: aTkn,
			body: api.CreateTransactionRequest{
				Amount:    6,
				AccountId: 1,
			},
			want: want{
				statusCode: http.StatusOK,
				account: api.Transaction{
					Id:       1001,
					OldSaldo: 436,
					NewSaldo: 430,
					Amount:   6,
					Created:  ptypes.TimestampNow(),
					Account: &api.Account{
						Id:          1,
						Name:        "Laverne Blackstock",
						Description: "Itchy Eye",
						Saldo:       430,
						NfcChipId:   "Hv8mnajqzIKO",
						Group: &api.Group{
							Id:          7,
							Name:        "PSS World Medical, Inc.",
							CanOverdraw: true,
						},
					},
				},
			},
		},
		{
			name:        "create new transaction with unkown account",
			accessToken: aTkn,
			body: api.CreateTransactionRequest{
				Amount:    6,
				AccountId: -45,
			},
			want: want{
				statusCode: http.StatusNotFound,
				errMsg:     "could not find account",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			body, err := json.Marshal(tt.body)
			is.NoErr(err) // could not marshal body

			req, err := http.NewRequest(http.MethodPost, RestUrlWithPath(fmt.Sprintf("v1/account/%d/transactions", tt.body.AccountId)), bytes.NewReader(body))
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

			var account api.Transaction
			err = jsonpb.Unmarshal(res.Body, &account)
			is.NoErr(err) // could not decode account

			checkTransactionEquality(t, account, tt.want.account)
		})
	}

}

func checkTransactionEquality(t *testing.T, got api.Transaction, want api.Transaction) {
	t.Helper()
	if got.Id != want.Id {
		t.Errorf("got id %d; wanted %d", got.Id, want.Id)
	}
	if got.OldSaldo != want.OldSaldo {
		t.Errorf("got oldsaldo %f; wanted %f", got.OldSaldo, want.OldSaldo)
	}
	if got.NewSaldo != want.NewSaldo {
		t.Errorf("got oldsaldo %f; wanted %f", got.NewSaldo, want.NewSaldo)
	}
	if got.Amount != want.Amount {
		t.Errorf("got amount %f; wanted %f", got.Amount, want.Amount)
	}
	if !reflect.DeepEqual(got.Account, want.Account) {
		t.Errorf("got account %v; wanted %v", got.Account, want.Account)
	}
	gotCreated, _ := ptypes.Timestamp(got.Created)
	wantCreated, _ := ptypes.Timestamp(want.Created)
	gotCreated = gotCreated.Round(time.Second)
	wantCreated = wantCreated.Round(time.Second)

	if !gotCreated.Equal(wantCreated) {
		t.Errorf("got created %s, wanted %s", gotCreated, wantCreated)
	}
}

func checkTotalCount(t *testing.T, transactionsResponse api.ListTransactionsResponse, total int32) {
	t.Helper()
	if transactionsResponse.TotalCount != total {
		t.Errorf("got totalcount %d, wanted %d", transactionsResponse.TotalCount, total)
	}
}

func checkTransactionLen(t *testing.T, transactionsResponse api.ListTransactionsResponse, length int) {
	t.Helper()
	if l := len(transactionsResponse.Transactions); l != length {
		t.Errorf("got %d transactions, wanted %d", l, length)
	}
}

func checkAccountId(t *testing.T, transactionsResponse api.ListTransactionsResponse, accountId int32) {
	t.Helper()
	for _, transaction := range transactionsResponse.Transactions {
		if transaction.Account.Id != accountId {
			t.Errorf("want only transactions for account with id %d, got one with account id %d", accountId, transaction.Account.Id)
		}
	}
}

func checkOrder(t *testing.T, transactionsResponse api.ListTransactionsResponse, order string) {
	t.Helper()
	if len(transactionsResponse.Transactions) < 2 {
		t.Fatalf("to few transactions to check order")
	}
	f := transactionsResponse.Transactions[0]
	s := transactionsResponse.Transactions[1]
	fCreated, err := ptypes.Timestamp(f.Created)
	if err != nil {
		t.Errorf("could not convert timestamp: %v", err)
	}
	sCreated, err := ptypes.Timestamp(s.Created)
	if err != nil {
		t.Errorf("could not convert timestamp: %v", err)
	}
	switch order {
	case "asc":
		if fCreated.After(sCreated) {
			t.Error("first transaction is created after the second, with order asc it should be before")
		}
	default:
		if fCreated.Before(sCreated) {
			t.Error("first transaction is created before the second, with order desc it should be after")
		}
	}
}
