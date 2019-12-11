package server

import (
	"context"
	"errors"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"google.golang.org/grpc"
)

var (
	ErrCouldNotCreateGroup = errors.New("could not create group")
)

type groupserver struct {
	storage models.GroupStorager
}

func RegisterGroupServer(s *grpc.Server, storage models.GroupStorager) {
	api.RegisterGroupsServiceServer(s, &groupserver{storage: storage})
}

func (g *groupserver) List(ctx context.Context, req *api.GroupListRequest) (*api.Groups, error) {
	groups, err := g.storage.GetAll()
	if err != nil {
		return nil, ErrGetAll
	}
	return groups, nil
}

func (g *groupserver) Create(ctx context.Context, req *api.GroupCreate) (*api.Group, error) {
	group, err := g.storage.Create(req.Name, req.Description, req.CanOverdraw)

	if err != nil {
		return nil, ErrCouldNotCreateGroup
	}

	return group, nil
}

func (g *groupserver) Get(ctx context.Context, req *api.IdRequest) (*api.Group, error) {
	group, err := g.storage.Read(req.Id)

	if err != nil {
		return nil, ErrNotFound
	}

	return group, nil
}

func (g *groupserver) Update(ctx context.Context, req *api.Group) (*api.Group, error) {
	group, err := g.storage.Update(req)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}
	return group, nil
}

func (g *groupserver) Delete(ctx context.Context, req *api.IdRequest) (*api.Status, error) {
	err := g.storage.Delete(req.Id)

	if err != nil {
		return &api.Status{
			Success:      false,
			ErrorMessage: ErrSomethingWentWrong.Error(),
		}, ErrSomethingWentWrong
	}

	return &api.Status{Success: true}, nil
}
