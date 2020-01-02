package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test"
)

func TestAccountserver_ListAccounts_Integration(t *testing.T) {
	test.IsIntegrationTest(t)
	teardown := startServers(t)
	defer teardown()

	atkn, _ := login(t)

	req, err := http.NewRequest(http.MethodGet, RestUrlWithPath("v1/accounts"), nil)
	req.Header.Add("Authorization", "Bearer "+atkn)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("could not request accounts: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		var jsonErr *api.Status
		err := json.NewDecoder(res.Body).Decode(&jsonErr)
		if err != nil {
			t.Fatalf("could not parse error: %v", err)
		}
		t.Logf("%v", jsonErr)
		return
	}

	var accounts api.ListAccountsResponse
	err = json.NewDecoder(res.Body).Decode(&accounts)

	if err != nil {
		t.Fatalf("could not parse accounts: %v", err)
	}

	if len(accounts.Accounts) != 100 {
		t.Errorf("got %d account, wanted %d", len(accounts.Accounts), 100)
	}
}
