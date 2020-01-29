package auth

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNoBearerAuth       = status.Error(codes.Unauthenticated, "bearer authorization required")
	ErrNoAuthHeader       = status.Error(codes.Unauthenticated, "authorization header required")
	ErrCouldNotAuthorize  = status.Error(codes.Internal, "authorization failed")
	ErrNoBasicAuth        = status.Error(codes.Unauthenticated, "basic authorization required")
	ErrNoUserNamePassword = status.Error(codes.Unauthenticated, "authorization required username and password")
)
