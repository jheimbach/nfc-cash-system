package handlers

import (
	"context"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/repositories"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type groupserver struct {
	storage repositories.GroupStorager
}

func RegisterGroupServer(s *grpc.Server, storage repositories.GroupStorager) {
	api.RegisterGroupsServiceServer(s, &groupserver{storage: storage})
}

func (g *groupserver) ListGroups(ctx context.Context, req *api.ListGroupsRequest) (*api.ListGroupsResponse, error) {
	var limit, offset int32
	if req.Paging != nil {
		limit = req.Paging.Limit
		offset = req.Paging.Offset
	}
	groups, count, err := g.storage.GetAll(ctx, limit, offset)

	if err != nil {
		return nil, ErrGetAll
	}
	return &api.ListGroupsResponse{
		Groups:     groups,
		TotalCount: int32(count),
	}, nil
}

func (g *groupserver) CreateGroup(ctx context.Context, req *api.CreateGroupRequest) (*api.Group, error) {
	group, err := g.storage.Create(ctx, req.Name, req.Description, req.CanOverdraw)

	if err != nil {
		return nil, ErrCouldNotCreateGroup
	}

	return group, nil
}

func (g *groupserver) GetGroup(ctx context.Context, req *api.GetGroupRequest) (*api.Group, error) {
	group, err := g.storage.Read(ctx, req.Id)

	if err != nil {
		return nil, ErrGroupNotFound
	}

	return group, nil
}

func (g *groupserver) UpdateGroup(ctx context.Context, req *api.Group) (*api.Group, error) {
	group, err := g.storage.Update(ctx, req)
	if err != nil {
		return nil, ErrSomethingWentWrong
	}
	return group, nil
}

func (g *groupserver) DeleteGroup(ctx context.Context, req *api.DeleteGroupRequest) (*empty.Empty, error) {
	err := g.storage.Delete(ctx, req.Id)

	if err != nil {
		if err == repositories.ErrNonEmptyDelete {
			return &empty.Empty{}, status.Error(codes.Aborted, "could not delete group, because it is not empty")
		}

		return &empty.Empty{}, ErrSomethingWentWrong
	}

	return &empty.Empty{}, nil
}
