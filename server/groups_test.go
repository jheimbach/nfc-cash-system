package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
)

type groupMockStorage struct {
	create func(name, description string, canOverdraw bool) (*api.Group, error)
	getAll func() (*api.Groups, error)
	read   func(id int32) (*api.Group, error)
	update func(group *api.Group) (*api.Group, error)
	delete func(id int32) error
}

func (g *groupMockStorage) Create(name, description string, canOverdraw bool) (*api.Group, error) {
	return g.create(name, description, canOverdraw)
}

func (g *groupMockStorage) GetAll() (*api.Groups, error) {
	return g.getAll()
}

func (g *groupMockStorage) Read(id int32) (*api.Group, error) {
	return g.read(id)
}

func (g *groupMockStorage) Update(group *api.Group) (*api.Group, error) {
	return g.update(group)
}

func (g *groupMockStorage) Delete(id int32) error {
	return g.delete(id)
}

func TestGroupserver_List(t *testing.T) {
	tests := []struct {
		name    string
		input   *api.GroupListRequest
		want    *api.Groups
		wantErr error
	}{
		{
			name:  "get simple list of accounts",
			input: &api.GroupListRequest{},
			want: &api.Groups{
				Groups: genGroupModels(2),
			},
		},
		{
			name:    "has error",
			input:   &api.GroupListRequest{},
			wantErr: ErrGetAll,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &groupserver{
				storage: &groupMockStorage{getAll: func() (*api.Groups, error) {
					if tt.wantErr != nil {
						return nil, sql.ErrNoRows
					}
					return tt.want, nil
				},
				},
			}
			got, err := a.List(context.Background(), tt.input)

			if err != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGroupserver_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   *api.Group
		want    *api.Group
		wantErr error
	}{
		{
			name: "create group",
			input: &api.Group{
				Name:        "test group",
				Description: "test",
				CanOverdraw: false,
			},
			want: &api.Group{
				Id:          1,
				Name:        "test group",
				Description: "test",
				CanOverdraw: false,
			},
		},
		{
			name: "create group",
			input: &api.Group{
				Name:        "test group",
				Description: "test",
				CanOverdraw: false,
			},
			wantErr: ErrCouldNotCreateGroup,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := groupserver{
				storage: &groupMockStorage{
					create: func(name, description string, canOverdraw bool) (*api.Group, error) {
						if tt.wantErr != nil {
							return nil, tt.wantErr
						}
						ret := tt.want
						ret.Id = 1
						return ret, nil
					},
				},
			}

			got, err := server.Create(context.Background(), tt.input)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err %v, expected %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("did not expect error, got %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestGroupserver_Update(t *testing.T) {
	tests := []struct {
		name      string
		want      *api.Group
		returnErr error // specifies the error which will be returned from storager
		wantErr   error
	}{
		{
			name: "update group",
			want: &api.Group{
				Id:          1,
				Name:        "testgroup",
				Description: "with description",
				CanOverdraw: false,
			},
			wantErr: nil,
		},
		{
			name:      "update group with id 0 returns error",
			returnErr: models.ErrModelNotSaved,
			wantErr:   ErrSomethingWentWrong,
		},
		{
			name:      "other error occured",
			returnErr: errors.New("some test error"),
			wantErr:   ErrSomethingWentWrong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := groupserver{
				storage: &groupMockStorage{
					update: func(group *api.Group) (*api.Group, error) {
						if tt.returnErr != nil {
							return nil, tt.returnErr
						}
						return group, nil
					},
				},
			}
			got, err := server.Update(context.Background(), tt.want)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err %v, expected %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error but got %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestGroupserver_Delete(t *testing.T) {
	tests := []struct {
		name    string
		request *api.IdRequest
		want    *api.Status
		wantErr error
	}{
		{
			name:    "delete group",
			request: &api.IdRequest{Id: 1},
			want:    &api.Status{Success: true},
		},
		{
			name:    "delete group with error",
			request: &api.IdRequest{Id: 1},
			want:    &api.Status{Success: false, ErrorMessage: ErrSomethingWentWrong.Error()},
			wantErr: ErrSomethingWentWrong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := groupserver{
				storage: &groupMockStorage{
					delete: func(id int32) error {
						if tt.wantErr != nil {
							return errors.New("test error")
						}
						return nil
					},
				},
			}

			got, err := server.Delete(context.Background(), tt.request)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err %v expected %v", err, tt.wantErr)
				}
			} else {
				if err != nil {
					t.Fatalf("got unexpected err %v", err)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, expected %v", got, tt.want)
			}
		})
	}
}

func genGroupModels(num int) []*api.Group {
	var groups []*api.Group
	for i := 0; i < num; i++ {
		groups = append(groups, &api.Group{
			Id:          int32(i + 1),
			Name:        fmt.Sprintf("group name %d", i+1),
			Description: "description",
			CanOverdraw: i%2 == 0,
		})
	}

	return groups
}
