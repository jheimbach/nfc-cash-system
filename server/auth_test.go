package server

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestRefreshTokenFromHeader(t *testing.T) {
	t.Run("extract refresh token from header", func(t *testing.T) {
		want := "<refresh_token>"
		md := metadata.New(map[string]string{"x-refresh-token": want})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		got, err := refreshTokenFromHeader(ctx)
		if err != nil {
			t.Errorf("did not expect error: %v", err)
		}

		if got != want {
			t.Errorf("got %s, wanted %s", got, want)
		}

	})

	t.Run("could not find refresh token in header", func(t *testing.T) {
		md := metadata.New(map[string]string{"refresh-token": ""})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := refreshTokenFromHeader(ctx)
		if err == nil {
			t.Errorf("expected an error")
		}

		want := status.Error(codes.Unauthenticated, "refresh token required")
		if !errors.Is(err, want) {
			t.Errorf("got err:%v, expected %v", err, want)
		}

	})
}

func TestAuthorizationHeader(t *testing.T) {
	t.Run("extract authorization from header", func(t *testing.T) {
		want := "<token>"
		md := metadata.New(map[string]string{"authorization": want})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		got, err := authorizationHeader(ctx)
		if err != nil {
			t.Errorf("did not expect error: %v", err)
		}

		if got != want {
			t.Errorf("got %s, wanted %s", got, want)
		}

	})

	t.Run("could not find authorization in header", func(t *testing.T) {
		md := metadata.New(map[string]string{"authentication": ""})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := authorizationHeader(ctx)
		if err == nil {
			t.Errorf("expected an error")
		}

		want := status.Error(codes.Unauthenticated, "authorization header required")
		if !errors.Is(err, want) {
			t.Errorf("got err:%v, expected %v", err, want)
		}

	})
}

func TestBearerAuthorization(t *testing.T) {
	t.Run("extract bearer token from header", func(t *testing.T) {
		want := "<token>"
		md := metadata.New(map[string]string{"authorization": "Bearer " + want})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		got, err := bearerAuthorization(ctx)
		if err != nil {
			t.Errorf("did not expect error: %v", err)
		}

		if got != want {
			t.Errorf("got %s, wanted %s", got, want)
		}

	})

	t.Run("could not find bearer token", func(t *testing.T) {
		md := metadata.New(map[string]string{"authorization": "Basic <token>"})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := bearerAuthorization(ctx)
		if err == nil {
			t.Errorf("expected an error")
		}

		want := status.Error(codes.Unauthenticated, "bearer authorization required")
		if !errors.Is(err, want) {
			t.Errorf("got err:%v, expected %v", err, want)
		}

	})
}

func TestBasicAuthorization(t *testing.T) {
	t.Run("extract basic token from header", func(t *testing.T) {
		want := []string{"user", "passwd"}
		md := metadata.New(map[string]string{"authorization": "Basic dXNlcjpwYXNzd2Q="})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		got, err := basicAuthorization(ctx)
		if err != nil {
			t.Errorf("did not expect error: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s, wanted %s", got, want)
		}

	})

	t.Run("could not find bearer token", func(t *testing.T) {
		md := metadata.New(map[string]string{"authorization": "Bearer <token>"})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := basicAuthorization(ctx)
		if err == nil {
			t.Errorf("expected an error")
		}

		want := status.Error(codes.Unauthenticated, "basic authorization required")
		if !errors.Is(err, want) {
			t.Errorf("got err:%v, expected %v", err, want)
		}

	})
}
