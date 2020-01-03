package handlers

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrGetAll                = status.Error(codes.NotFound, "could not load list of accounts")
	ErrCouldNotCreateAccount = status.Error(codes.Internal, "could not save new account")
	ErrNotFound              = status.Error(codes.NotFound, "could not find account")
	ErrSomethingWentWrong    = status.Error(codes.Internal, "something went wrong")
	ErrNameOrPasswdWrong     = status.Error(codes.Unauthenticated, "username or password wrong")
	ErrNoRefreshToken        = status.Error(codes.Unauthenticated, "refresh token required")
	ErrCouldNotLogOut        = status.Error(codes.Internal, "could not log user out")
	ErrCouldNotCreateGroup   = status.Error(codes.Internal, "could not create group")
)
