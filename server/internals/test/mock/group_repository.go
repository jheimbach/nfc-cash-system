package mock

import (
	"context"

	"github.com/JHeimbach/nfc-cash-system/server/api"
)

type GroupRepository struct {
	GetAllByIdsFunc func(ids []int32) (map[int32]*api.Group, error)
	CreateFunc      func(string, string, bool) (*api.Group, error)
	GetAllFunc      func(int32, int32) ([]*api.Group, int, error)
	ReadFunc        func(int32) (*api.Group, error)
	UpdateFunc      func(*api.Group) (*api.Group, error)
	DeleteFunc      func(int32) error
}

func (g *GroupRepository) GetAllByIds(_ context.Context, ids []int32) (map[int32]*api.Group, error) {
	return g.GetAllByIdsFunc(ids)
}

func (g *GroupRepository) Create(_ context.Context, name, desc string, overdraw bool) (*api.Group, error) {
	return g.CreateFunc(name, desc, overdraw)
}

func (g *GroupRepository) GetAll(_ context.Context, limit, offset int32) ([]*api.Group, int, error) {
	return g.GetAllFunc(limit, offset)
}

func (g *GroupRepository) Read(_ context.Context, id int32) (*api.Group, error) {
	return g.ReadFunc(id)
}

func (g *GroupRepository) Update(_ context.Context, group *api.Group) (*api.Group, error) {
	return g.UpdateFunc(group)
}

func (g *GroupRepository) Delete(_ context.Context, id int32) error {
	return g.DeleteFunc(id)
}
