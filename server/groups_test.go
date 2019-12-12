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
	"github.com/golang/protobuf/ptypes/empty"
)

type groupMockStorage struct {
	create      func(name, description string, canOverdraw bool) (*api.Group, error)
	getAll      func() ([]*api.Group, error)
	read        func(id int32) (*api.Group, error)
	update      func(group *api.Group) (*api.Group, error)
	delete      func(id int32) error
	getAllByIds func(ids []int32) (map[int32]*api.Group, error)
}

func (g *groupMockStorage) Create(name, description string, canOverdraw bool) (*api.Group, error) {
	return g.create(name, description, canOverdraw)
}

func (g *groupMockStorage) GetAll() ([]*api.Group, error) {
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
func (g groupMockStorage) GetAllByIds(ids []int32) (map[int32]*api.Group, error) {
	return g.getAllByIds(ids)
}

func TestGroupserver_List(t *testing.T) {
	tests := []struct {
		name    string
		input   *api.ListGroupsRequest
		want    *api.ListGroupsResponse
		wantErr error
	}{
		{
			name:  "get simple list of accounts",
			input: &api.ListGroupsRequest{},
			want: &api.ListGroupsResponse{
				Groups:     genGroupModels(2),
				TotalCount: 2,
			},
		},
		{
			name:    "has error",
			input:   &api.ListGroupsRequest{},
			wantErr: ErrGetAll,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &groupserver{
				storage: &groupMockStorage{getAll: func() ([]*api.Group, error) {
					if tt.wantErr != nil {
						return nil, sql.ErrNoRows
					}
					return tt.want.Groups, nil
				},
				},
			}
			got, err := a.ListGroups(context.Background(), tt.input)

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

/*

groups:<id:1 name:"group name 1" description:"description" can_overdraw:true > groups:<id:2 name:"group name 2" description:"description" > total_count:2
groups:<id:1 name:"group name 1" description:"description" can_overdraw:true > groups:<id:2 name:"group name 2" description:"description" >

*/

func TestGroupserver_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   *api.CreateGroupRequest
		want    *api.Group
		wantErr error
	}{
		{
			name: "create group",
			input: &api.CreateGroupRequest{
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
			input: &api.CreateGroupRequest{
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

			got, err := server.CreateGroup(context.Background(), tt.input)
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
			got, err := server.UpdateGroup(context.Background(), tt.want)

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
		request *api.DeleteGroupRequest
		want    *empty.Empty
		wantErr error
	}{
		{
			name:    "delete group",
			request: &api.DeleteGroupRequest{Id: 1},
			want:    &empty.Empty{},
		},
		{
			name:    "delete group with error",
			request: &api.DeleteGroupRequest{Id: 1},
			want:    &empty.Empty{},
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

			got, err := server.DeleteGroup(context.Background(), tt.request)

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
