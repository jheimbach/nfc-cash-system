package handlers

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test"
)

func TestAccountserver_ListAccounts_Integration(t *testing.T) {
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
		AccessToken  string
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
			AccessToken: aTkn,
			want: want{
				statusCode:   http.StatusOK,
				accountsLen:  100,
				accountTotal: 100,
			},
		},
		{
			name:        "get first 10 accounts",
			AccessToken: aTkn,
			pagingLimit: 10,
			want: want{
				statusCode:   http.StatusOK,
				accountsLen:  10,
				accountTotal: 100,
			},
		},
		{
			name:         "get second 10 accounts",
			AccessToken:  aTkn,
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
			AccessToken: aTkn,
			groupId:     1,
			want: want{
				statusCode:   http.StatusOK,
				accountsLen:  10,
				accountTotal: 10,
			},
		},
		{
			name:        "filter by group with limit",
			AccessToken: aTkn,
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
			AccessToken:  aTkn,
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
			if tt.AccessToken != "" {
				req.Header.Add("Authorization", "Bearer "+tt.AccessToken)
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("could not request accounts: %v", err)
			}
			defer res.Body.Close()

			if tt.want.statusCode != http.StatusOK {
				var jsonErr *api.Status
				err := json.NewDecoder(res.Body).Decode(&jsonErr)
				if err != nil {
					t.Fatalf("could not parse error: %v", err)
				}

				if jsonErr.Message != tt.want.errMsg {
					t.Errorf("got err msg: %q, wanted: %q", jsonErr.Message, tt.want.errMsg)
				}
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
