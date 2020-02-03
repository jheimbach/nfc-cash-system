package auth

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/jheimbach/nfc-cash-system/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var bypassAuth = map[string]struct{}{
	"/api.UserService/AuthenticateUser": {},
	"/api.HealthService/Health":         {},
	"/api.UserService/RefreshToken":     {},
}

func InitInterceptor(gen TokenGenerator) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if _, ok := bypassAuth[info.FullMethod]; ok {
			return handler(ctx, req)
		}
		token, err := bearerAuthorization(ctx)
		if err != nil {
			return nil, err
		}
		user, _, err := gen.VerifyToken(token, AccessToken)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, "user", user)

		return handler(ctx, req)
	}
}

func bearerAuthorization(ctx context.Context) (string, error) {
	header, err := authorizationHeader(ctx)
	if err != nil {
		return "", err
	}

	// check if authorization is a basic auth
	authorization := strings.SplitN(header, " ", 2)
	if len(authorization) != 2 || authorization[0] != "Bearer" {
		return "", ErrNoBearerAuth
	}

	return authorization[1], nil
}

func authorizationHeader(ctx context.Context) (string, error) {
	// load metadata
	mb, _ := metadata.FromIncomingContext(ctx)

	// get authorization metadata (header)
	authHeader := mb.Get("authorization")
	if len(authHeader) < 1 {
		return "", ErrNoAuthHeader
	}
	return authHeader[0], nil
}

func UsernameAndPasswortFromContext(ctx context.Context) ([]string, error) {
	header, err := authorizationHeader(ctx)
	if err != nil {
		return nil, err
	}
	// check if authorization is a basic auth
	authorization := strings.SplitN(header, " ", 2)
	if len(authorization) != 2 || authorization[0] != "Basic" {
		return nil, ErrNoBasicAuth
	}

	// decode username and password
	payload, _ := base64.StdEncoding.DecodeString(authorization[1])
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return nil, ErrNoUserNamePassword
	}

	return pair, nil
}

func RetrieveUserFromContext(ctx context.Context) (*api.User, error) {
	user, ok := ctx.Value("user").(*api.User)
	if !ok || user == nil {
		return nil, ErrCouldNotAuthorize
	}
	return user, nil
}
