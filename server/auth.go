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

func RegisterUserServer(s *grpc.Server, storage models.UserStorager) {
	api.RegisterUserServiceServer(s, &userServer{
		storage: storage,
	})
}

func (a *userServer) AuthenticateUser(ctx context.Context, empty *empty.Empty) (*api.AuthenticateResponse, error) {
	pair, err := basicAuthorization(ctx)
	if err != nil {
		return nil, err
	}

	// authenticate user from database
	user, err := a.storage.Authenticate(ctx, pair[0], pair[1])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "username or password wrong")
	}

	// create access token
	expire := auth.ExpirationTime(5 * time.Minute)
	accessToken, err := auth.CreateAccessToken(user, expire)
	if err != nil {
		return nil, status.Error(codes.Internal, "authorization failed")
	}

	// create refresh token
	refreshToken, err := a.refreshToken(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "authorization failed")
	}

	return &api.AuthenticateResponse{
		TokenType:    api.AuthenticateResponse_BEARER,
		AccessToken:  accessToken,
		ExpiresIn:    expire.UnixNano(),
		RefreshToken: refreshToken,
	}, nil
}

func (a *userServer) LogoutUser(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	recvToken, err := bearerAuthorization(ctx)
	if err != nil {
		return nil, err
	}

	user, err := auth.VerifyToken(recvToken, auth.AccessTokenKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "authentication failed")
	}

	err = a.storage.DeleteRefreshKey(ctx, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "could not log user out")
	}

	return &empty.Empty{}, nil
}

func (a *userServer) RefreshToken(ctx context.Context, e *empty.Empty) (*api.AuthenticateResponse, error) {
	aToken, err := bearerAuthorization(ctx)
	if err != nil {
		return nil, err
	}
	rToken, err := refreshTokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}
	user, err := auth.VerifyToken(aToken, auth.AccessTokenKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "authentication failed")
	}
	refreshK, err := a.storage.GetRefreshKey(ctx, user.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "authentication failed")
	}
	_, err = auth.VerifyToken(rToken, refreshK)
	if err != nil {
		return nil, status.Error(codes.Internal, "authentication failed")
	}

	expire := auth.ExpirationTime(5 * time.Minute)
	newAToken, err := auth.CreateAccessToken(user, expire)
	if err != nil {
		return nil, status.Error(codes.Internal, "authentication failed")
	}

	return &api.AuthenticateResponse{
		TokenType:    api.AuthenticateResponse_BEARER,
		AccessToken:  newAToken,
		RefreshToken: rToken,
		ExpiresIn:    expire.Unix(),
	}, nil
}

func (a *userServer) refreshToken(ctx context.Context, user *api.User) (string, error) {
	// create random key for each user
	refreshKey := auth.CreateRandomKey()

	// create jwt token with userId and random key
	refreshToken, err := auth.CreateToken(user, auth.ExpirationTime(7*24*time.Hour), refreshKey)
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

func basicAuthorization(ctx context.Context) ([]string, error) {
	header, err := authorizationHeader(ctx)
	if err != nil {
		return nil, err
	}
	// check if authorization is a basic auth
	authorization := strings.SplitN(header, " ", 2)
	if len(authorization) != 2 || authorization[0] != "Basic" {
		return nil, status.Error(codes.Unauthenticated, "basic authorization required")
	}

	// decode username and password
	payload, _ := base64.StdEncoding.DecodeString(authorization[1])
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return nil, status.Error(codes.Unauthenticated, "authorization required username and password")
	}

	return pair, nil
}

func bearerAuthorization(ctx context.Context) (string, error) {
	header, err := authorizationHeader(ctx)
	if err != nil {
		return "", err
	}

	// check if authorization is a basic auth
	authorization := strings.SplitN(header, " ", 2)
	if len(authorization) != 2 || authorization[0] != "Bearer" {
		return "", status.Error(codes.Unauthenticated, "bearer authorization required")
	}

	return authorization[1], nil
}

func authorizationHeader(ctx context.Context) (string, error) {
	// load metadata
	mb, _ := metadata.FromIncomingContext(ctx)

	// get authorization metadata (header)
	authHeader := mb.Get("authorization")
	if len(authHeader) < 1 {
		return "", status.Error(codes.Unauthenticated, "authorization header required")
	}
	return authHeader[0], nil
}

func refreshTokenFromHeader(ctx context.Context) (string, error) {
	// load metadata
	mb, _ := metadata.FromIncomingContext(ctx)

	// get authorization metadata (header)
	token := mb.Get("x-refresh-token")
	if len(token) < 1 {
		return "", status.Error(codes.Unauthenticated, "refresh token required")
	}
	return token[0], nil
}
