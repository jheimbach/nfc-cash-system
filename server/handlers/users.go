package handlers

import (
	"context"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/auth"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type userServer struct {
	storage        models.UserStorager
	tokenGenerator auth.TokenGenerator
}

func RegisterUserServer(s *grpc.Server, storage models.UserStorager) {
	api.RegisterUserServiceServer(s, &userServer{
		storage:        storage,
		tokenGenerator: auth.NewJwtGenerator(),
	})
}

func (a *userServer) AuthenticateUser(ctx context.Context, empty *empty.Empty) (*api.AuthenticateResponse, error) {
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
	accessToken, err := a.tokenGenerator.CreateToken(user, expire, auth.AccessTokenKey)
	if err != nil {
		return nil, auth.ErrCouldNotAuthorize
	}

	// create refresh token
	refreshToken, err := a.createRefreshToken(ctx, user)
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
	user, err := auth.RetrieveUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = a.storage.DeleteRefreshKey(ctx, user.Id)
	if err != nil {
		return nil, ErrCouldNotLogOut
	}

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

	refreshK, err := a.storage.GetRefreshKey(ctx, user.Id)
	if err != nil {
		return nil, auth.ErrCouldNotAuthorize
	}
	_, err = a.tokenGenerator.VerifyToken(rToken, refreshK)
	if err != nil {
		return nil, auth.ErrCouldNotAuthorize
	}

	expire := a.tokenGenerator.ExpirationTime(5 * time.Minute)
	newAToken, err := a.tokenGenerator.CreateToken(user, expire, auth.AccessTokenKey)
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

func (a *userServer) createRefreshToken(ctx context.Context, user *api.User) (string, error) {
	// create random key for each user
	refreshKey := a.tokenGenerator.CreateRandomKey()

	// create jwt token with userId and random key
	refreshToken, err := a.tokenGenerator.CreateToken(user, a.tokenGenerator.ExpirationTime(7*24*time.Hour), refreshKey)
	if err != nil {
		return "", err
	}

	//save refresh key to database
	err = a.storage.InsertRefreshKey(ctx, user.Id, refreshKey)
	if err != nil {
		if err == models.ErrUserHasRefreshKey {
			dErr := a.storage.DeleteRefreshKey(ctx, user.Id)
			if dErr != nil {
				return "", dErr
			}
		}
		if err == models.ErrRefreshKeyIsInUse || err == models.ErrUserHasRefreshKey {
			return a.createRefreshToken(ctx, user)
		}
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
