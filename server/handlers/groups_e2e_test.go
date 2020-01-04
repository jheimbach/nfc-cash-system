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

func TestGroupserver_E2E_ListGroups(t *testing.T) {
	test.IsIntegrationTest(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode int
		errMsg     string
		groupsLen  int
		groupTotal int32
	}
	tests := []struct {
		name         string
		accessToken  string
		want         want
		pagingLimit  int
		pagingOffset int
	}{
		{
			name: "no accesstoken given",
			want: want{
				statusCode: http.StatusUnauthorized,
				errMsg:     "authorization header required",
			},
		},
		{
			name:        "get all groups",
			accessToken: aTkn,
			want: want{
				statusCode: http.StatusOK,
				groupsLen:  10,
				groupTotal: 10,
			},
		},
		{
			name:        "get groups with limit",
			accessToken: aTkn,
			pagingLimit: 5,
			want: want{
				statusCode: http.StatusOK,
				groupsLen:  5,
				groupTotal: 10,
			},
		},
		{
			name:         "get groups with limit and offset",
			accessToken:  aTkn,
			pagingLimit:  5,
			pagingOffset: 8,
			want: want{
				statusCode: http.StatusOK,
				groupsLen:  2,
				groupTotal: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := RestUrlWithPath("v1/groups")
			path, err := url.Parse(p)
			if err != nil {
				t.Fatalf("could not parse url: %s; %v", p, err)
			}

			if tt.pagingLimit != 0 {
				q := path.Query()
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
				t.Fatalf("could not request groups: %v", err)
			}
			defer res.Body.Close()

			if tt.want.statusCode != http.StatusOK {
				checkError(t, res, tt.want.statusCode, tt.want.errMsg)
				return
			}

			var groups api.ListGroupsResponse
			err = json.NewDecoder(res.Body).Decode(&groups)

			if err != nil {
				t.Fatalf("could not parse groups: %v", err)
			}

			if l := len(groups.Groups); l != tt.want.groupsLen {
				t.Errorf("got %d groups, wanted %d", l, tt.want.groupsLen)
			}

			if groups.TotalCount != tt.want.groupTotal {
				t.Errorf("got totalcount %d, wanted %d", groups.TotalCount, tt.want.groupTotal)
			}
		})
	}
}

func TestGroupserver_E2E_CreateGroup(t *testing.T) {
	is := isPkg.New(t)
	test.IsIntegrationTest(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode int
		errMsg     string
		group      api.Group
	}
	tests := []struct {
		name        string
		accessToken string
		want        want
		body        *api.CreateGroupRequest
	}{
		{
			name: "no accesstoken given",
			want: want{
				statusCode: http.StatusUnauthorized,
				errMsg:     "authorization header required",
			},
		},
		{
			name:        "create minimal group",
			accessToken: aTkn,
			want: want{
				statusCode: http.StatusOK,
				group: api.Group{
					Id:   11,
					Name: "testgroup",
				},
			},
			body: &api.CreateGroupRequest{
				Name: "testgroup",
			},
		},
		{
			name:        "create group with description",
			accessToken: aTkn,
			want: want{
				statusCode: http.StatusOK,
				group: api.Group{
					Id:          12,
					Name:        "testgroup",
					Description: "desc",
				},
			},
			body: &api.CreateGroupRequest{
				Name:        "testgroup",
				Description: "desc",
			},
		}, {
			name:        "create group with can overdraw",
			accessToken: aTkn,
			want: want{
				statusCode: http.StatusOK,
				group: api.Group{
					Id:          13,
					Name:        "testgroup",
					CanOverdraw: true,
				},
			},
			body: &api.CreateGroupRequest{
				Name:        "testgroup",
				CanOverdraw: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			body, err := json.Marshal(tt.body)
			is.NoErr(err) // could not marshal body

			req, err := http.NewRequest(http.MethodPost, RestUrlWithPath("v1/groups"), bytes.NewReader(body))
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

			var group api.Group
			err = json.NewDecoder(res.Body).Decode(&group)
			is.NoErr(err) // could not decode group

			is.Equal(group, tt.want.group) // group is not the expected
		})
	}
}

func TestGroupserver_E2E_GetGroup(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode int
		errMsg     string
		group      api.Group
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
			name:        "get group with id 1",
			accessToken: aTkn,
			accountId:   1,
			want: want{
				statusCode: http.StatusOK,
				group: api.Group{
					Id:   1,
					Name: "H2O Plus",
				},
			},
		},
		{
			name:        "get group that does not exist",
			accessToken: aTkn,
			accountId:   -45,
			want: want{
				statusCode: http.StatusNotFound,
				errMsg:     "could not find group",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			req, err := http.NewRequest(http.MethodGet, RestUrlWithPath(fmt.Sprintf("v1/group/%d", tt.accountId)), nil)
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

			var group api.Group
			err = json.NewDecoder(res.Body).Decode(&group)
			is.NoErr(err) // could not decode group

			is.Equal(group, tt.want.group) // group is not the expected
		})
	}
}

func TestGroupserver_E2E_UpdateGroup(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	type want struct {
		statusCode int
		errMsg     string
		group      api.Group
	}
	tests := []struct {
		name        string
		accessToken string
		body        *api.Group
		want        want
	}{
		{
			name: "no accesstoken given",
			body: &api.Group{
				Id: 1,
			},
			want: want{
				statusCode: http.StatusUnauthorized,
				errMsg:     "authorization header required",
			},
		},
		{
			name:        "update group name",
			accessToken: aTkn,
			body: &api.Group{
				Id:   1,
				Name: "H20 Plus",
			},
			want: want{
				statusCode: http.StatusOK,
				group: api.Group{
					Id:   1,
					Name: "H20 Plus",
				},
			},
		},

		{
			name:        "update group description",
			accessToken: aTkn,
			body: &api.Group{
				Id:          1,
				Name:        "H20 Plus",
				Description: "test",
			},
			want: want{
				statusCode: http.StatusOK,
				group: api.Group{
					Id:          1,
					Name:        "H20 Plus",
					Description: "test",
				},
			},
		},
		{
			name:        "update group can overdraw",
			accessToken: aTkn,
			body: &api.Group{
				Id:          1,
				Name:        "H20 Plus",
				Description: "test",
				CanOverdraw: true,
			},
			want: want{
				statusCode: http.StatusOK,
				group: api.Group{
					Id:          1,
					Name:        "H20 Plus",
					Description: "test",
					CanOverdraw: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			body, err := json.Marshal(tt.body)
			is.NoErr(err) // could not marshal body

			req, err := http.NewRequest(http.MethodPut, RestUrlWithPath(fmt.Sprintf("v1/group/%d", tt.body.Id)), bytes.NewReader(body))
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

			var group api.Group
			err = json.NewDecoder(res.Body).Decode(&group)
			is.NoErr(err) // could not decode group

			is.Equal(group, tt.want.group) // group is not the expected
		})
	}
}

func TestGroupserver_E2E_DeleteGroup(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	teardown := startServers(t)
	defer teardown()

	aTkn, _ := login(t)

	tests := []struct {
		name             string
		accessToken      string
		groupId          int
		statusCode       int
		createEmptyGroup bool
		errMsg           string
	}{
		{
			name:       "no accesstoken given",
			groupId:    1,
			statusCode: http.StatusUnauthorized,
			errMsg:     "authorization header required",
		},
		{
			name:             "delete empty group",
			accessToken:      aTkn,
			createEmptyGroup: true,
			statusCode:       http.StatusOK,
		},
		{
			name:        "delete non empty group",
			accessToken: aTkn,
			groupId:     1,
			statusCode:  http.StatusConflict,
			errMsg:      "could not delete group, because it is not empty",
		},
		{
			name:        "delete group with invalid id",
			accessToken: aTkn,
			groupId:     -45,
			statusCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			groupId := tt.groupId
			if tt.createEmptyGroup {
				groupId = createGroupAndReturnId(t)
			}
			req, err := http.NewRequest(http.MethodDelete, RestUrlWithPath(fmt.Sprintf("v1/group/%d", groupId)), nil)
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
			if res.StatusCode != http.StatusOK {
				t.Errorf("got statuscode %d, expected 200", res.StatusCode)
			}

			b, err := ioutil.ReadAll(res.Body)
			is.NoErr(err) // could not read body

			if string(b) != "{}" {
				t.Errorf("expected empty body, got %q", b)
			}
		})
	}
}

func createGroupAndReturnId(t *testing.T) int {
	is := isPkg.New(t)
	body, err := json.Marshal(api.CreateGroupRequest{
		Name: "testgroup",
	})
	is.NoErr(err) //could not marshal empty group request
	req, err := http.NewRequest(http.MethodPost, RestUrlWithPath("v1/groups"), bytes.NewReader(body))
	is.NoErr(err) //could not create post request

	res, err := http.DefaultClient.Do(req)
	is.NoErr(err) //could not create post request

	var group api.Group
	err = json.NewDecoder(res.Body).Decode(&group)
	is.NoErr(err) // could not decode group

	err = res.Body.Close()
	is.NoErr(err) // could not close result body

	return int(group.Id)
}
