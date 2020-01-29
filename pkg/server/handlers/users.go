package handlers

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jheimbach/nfc-cash-system/api"
	"github.com/jheimbach/nfc-cash-system/pkg/server/auth"
	"github.com/jheimbach/nfc-cash-system/pkg/server/repositories"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type userServer struct {
	storage        repositories.UserStorager
	tokenGenerator auth.TokenGenerator
}

func RegisterUserServer(s *grpc.Server, storage repositories.UserStorager, generator auth.TokenGenerator) {
	api.RegisterUserServiceServer(s, &userServer{
		storage:        storage,
		tokenGenerator: generator,
	})
}

func (a *userServer) AuthenticateUser(ctx context.Context, _ *empty.Empty) (*api.AuthenticateResponse, error) {
	pair, err := auth.UsernameAndPasswortFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// authenticate user from database
	user, err := a.storage.Authenticate(ctx, pair[0], pair[1])
	if err != nil {
		return nil, ErrNameOrPasswdWrong
	}

	// create access token
	expire := a.tokenGenerator.ExpirationTime(5 * time.Minute)
	accessToken, err := a.tokenGenerator.CreateToken(user, expire, auth.AccessToken)
	if err != nil {
		return nil, auth.ErrCouldNotAuthorize
	}

	// create refresh token
	refreshToken, err := a.createRefreshToken(user)
	if err != nil {
		return nil, auth.ErrCouldNotAuthorize
	}

	return &api.AuthenticateResponse{
		TokenType:    api.AuthenticateResponse_BEARER,
		AccessToken:  accessToken,
		ExpiresIn:    expire.Unix(),
		RefreshToken: refreshToken,
	}, nil
}

func (a *userServer) LogoutUser(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	// todo: set refreshtoken and accesstoken on a blacklist
	return &empty.Empty{}, nil
}

func (a *userServer) RefreshToken(ctx context.Context, e *empty.Empty) (*api.AuthenticateResponse, error) {
	user, err := auth.RetrieveUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	rToken, err := refreshTokenFromHeader(ctx)
	if err != nil {
		return nil, err
	}

	_, expires, err := a.tokenGenerator.VerifyToken(rToken, auth.RefreshToken)
	if err != nil {
		return nil, auth.ErrCouldNotAuthorize
	}

	// if refreshToken is about to expire, renew it
	if expires.After(time.Now().Add(-10 * time.Minute)) {
		rToken, err = a.createRefreshToken(user)
		if err != nil {
			return nil, err
		}
	}

	expire := a.tokenGenerator.ExpirationTime(5 * time.Minute)
	newAToken, err := a.tokenGenerator.CreateToken(user, expire, auth.AccessToken)
	if err != nil {
		return nil, auth.ErrCouldNotAuthorize
	}

	return &api.AuthenticateResponse{
		TokenType:    api.AuthenticateResponse_BEARER,
		AccessToken:  newAToken,
		RefreshToken: rToken,
		ExpiresIn:    expire.Unix(),
	}, nil
}

func (a *userServer) createRefreshToken(user *api.User) (string, error) {

	// create jwt token with userId and random key
	refreshToken, err := a.tokenGenerator.CreateToken(user, a.tokenGenerator.ExpirationTime(time.Hour), auth.RefreshToken)
	if err != nil {
		return "", err
	}

	// return jwt refresh token
	return refreshToken, nil
}

func refreshTokenFromHeader(ctx context.Context) (string, error) {
	// load metadata
	mb, _ := metadata.FromIncomingContext(ctx)

	// get authorization metadata (header)
	token := mb.Get("x-refresh-token")
	if len(token) < 1 {
		return "", ErrNoRefreshToken
	}
	return token[0], nil
}
