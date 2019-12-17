package server

import (
	"context"
	"encoding/base64"
	"strings"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/auth"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type userServer struct {
	storage models.UserStorager
	jwtKey  string
}

func RegisterAuthServer(s *grpc.Server, storage models.UserStorager) {
	api.RegisterUserServiceServer(s, &userServer{
		storage: storage,
	})
}

func (a *userServer) AuthenticateUser(ctx context.Context, empty *empty.Empty) (*api.AuthenticateResponse, error) {
	// load metadata
	mb, _ := metadata.FromIncomingContext(ctx)

	// get authorization metadata (header)
	authHeader := mb.Get("authorization")
	if len(authHeader) < 1 {
		return nil, status.Error(codes.Unauthenticated, "authorization header required")
	}

	// check if authorization is a basic auth
	authorization := strings.SplitN(authHeader[0], " ", 2)
	if len(authorization) != 2 || authorization[0] != "Basic" {
		return nil, status.Error(codes.Unauthenticated, "basic authorization required")
	}

	// decode username and password
	payload, _ := base64.StdEncoding.DecodeString(authorization[1])
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return nil, status.Error(codes.Unauthenticated, "authorization required username and password")
	}

	// authenticate user from database
	user, err := a.storage.Authenticate(ctx, pair[0], pair[1])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "username or password wrong")
	}

	// create access token
	expire := auth.ExpirationTime(5 * time.Minute)
	accessToken, err := auth.CreateAccessToken(user, expire, auth.AccessTokenKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "authorization failed")
	}

	// create refresh token
	refreshToken, err := a.refreshToken(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "authorization failed")
	}

	// return response
	return &api.AuthenticateResponse{
		TokenType:    api.AuthenticateResponse_BEARER,
		AccessToken:  accessToken,
		ExpiresIn:    expire.UnixNano(),
		RefreshToken: refreshToken,
	}, nil
}

func (a *userServer) LogoutUser(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (a *userServer) refreshToken(ctx context.Context, user *api.User) (string, error) {
	// create random key for each user
	refreshKey := auth.CreateRandomKey()

	// create jwt token with userId and random key
	refreshToken, err := auth.CreateRefreshToken(user, auth.ExpirationTime(7*24*time.Hour), refreshKey)
	if err != nil {
		return "", err
	}

	//save refresh key to database
	err = a.storage.InsertRefreshKey(ctx, user.Id, refreshKey)
	if err != nil {
		if err == models.ErrUserHasRefreshKey || err == models.ErrRefreshKeyIsInUse {
			err = a.storage.DeleteRefreshKey(ctx, user.Id)
			if err != nil {
				return "", err
			}
			return a.refreshToken(ctx, user)
		}
		return "", err
	}

	// return jwt refresh token
	return refreshToken, nil
}
