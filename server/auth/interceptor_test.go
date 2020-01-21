package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

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

		want := ErrNoAuthHeader
		if err != want {
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

		want := ErrNoBearerAuth
		if err != want {
			t.Errorf("got err:%v, expected %v", err, want)
		}

	})

	t.Run("could not find authorization header token", func(t *testing.T) {
		md := metadata.New(map[string]string{})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := bearerAuthorization(ctx)
		if err == nil {
			t.Errorf("expected an error")
		}

		want := ErrNoAuthHeader
		if err != want {
			t.Errorf("got err:%v, expected %v", err, want)
		}
	})
}

func TestBasicAuthorization(t *testing.T) {
	t.Run("extract basic token from header", func(t *testing.T) {
		want := []string{"user", "passwd"}
		md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("user:passwd")))})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		got, err := UsernameAndPasswortFromContext(ctx)
		if err != nil {
			t.Errorf("did not expect error: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %s, wanted %s", got, want)
		}

	})

	t.Run("returns error if no basic authorization is found", func(t *testing.T) {
		md := metadata.New(map[string]string{"authorization": "Bearer <token>"})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := UsernameAndPasswortFromContext(ctx)
		if err == nil {
			t.Errorf("expected an error")
		}

		want := ErrNoBasicAuth
		if err != want {
			t.Errorf("got err:%v, expected %v", err, want)
		}

	})

	t.Run("returns error if no header is found", func(t *testing.T) {
		md := metadata.New(map[string]string{})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := UsernameAndPasswortFromContext(ctx)
		if err == nil {
			t.Errorf("expected an error")
		}

		want := ErrNoAuthHeader
		if err != want {
			t.Errorf("got err:%v, expected %v", err, want)
		}

	})
	t.Run("returns error if is could not find username:password", func(t *testing.T) {
		md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("usernamepassword")))})
		ctx := metadata.NewIncomingContext(context.Background(), md)

		_, err := UsernameAndPasswortFromContext(ctx)

		want := ErrNoUserNamePassword
		if err != want {
			t.Errorf("got err:%v, expected %v", err, want)
		}
	})
}

func TestUnaryInterceptor(t *testing.T) {
	mUser := &api.User{
		Id:      1,
		Name:    "testuser1",
		Email:   "test@example.com",
		Created: mockTimeStamp(),
	}
	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}
	type gen struct {
		user *api.User
		key  TokenType
		exp  time.Time
	}
	tests := []struct {
		name     string
		info     *grpc.UnaryServerInfo
		tokenGen gen
		header   map[string]string
		wantErr  error
		handler  func(ctx context.Context, req interface{}) (interface{}, error)
	}{
		{
			name: "valid token",
			info: &grpc.UnaryServerInfo{FullMethod: "test/url"},
			tokenGen: gen{
				user: mUser,
				key:  AccessToken,
				exp:  time.Now().Add(5 * time.Minute),
			},
			header: map[string]string{},
			handler: func(ctx context.Context, req interface{}) (i interface{}, err error) {
				user, ok := ctx.Value("user").(*api.User)
				if !ok {
					t.Errorf("could not find user in request context")
				}

				if !reflect.DeepEqual(user, mUser) {
					t.Errorf("got wrong user from request")
				}
				return nil, nil
			},
		},
		{
			name:    "invalid token",
			info:    &grpc.UnaryServerInfo{FullMethod: "test/url"},
			header:  map[string]string{"authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoxLCJuYW1lIjoidGVzdHVzZXIxIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiY3JlYXRlZCI6eyJzZWNvbmRzIjoxNTQ3ODMxNzc1fX0sImV4cCI6NTE3Njg1OTQ3OCwianRpIjoiNmQ0NGFlNmY3MzBiMmJhZDMwYmY3MGM3NzQ1Nzc2YmIiLCJzdWIiOiJ1c2VyX3Rlc3R1c2VyMV8xIn0.y5r3D4NuQa53gHAk79HF1N9OUuRhWwNF5Dj-vXxWMgY"},
			wantErr: ErrCouldNotAuthorize,
			handler: mockHandler,
		},
		{
			name:    "no token send",
			info:    &grpc.UnaryServerInfo{FullMethod: "test/url"},
			header:  map[string]string{"authorization": ""},
			wantErr: ErrNoBearerAuth,
			handler: mockHandler,
		},
		{
			name:    "no auth header send",
			info:    &grpc.UnaryServerInfo{FullMethod: "test/url"},
			header:  map[string]string{},
			wantErr: ErrNoAuthHeader,
			handler: mockHandler,
		},
		{
			name:   "is AuthenticateUser route",
			info:   &grpc.UnaryServerInfo{FullMethod: "/api.UserService/AuthenticateUser"},
			header: map[string]string{"authorization": ""},
			handler: func(ctx context.Context, req interface{}) (i interface{}, err error) {
				_, ok := ctx.Value("user").(*api.User)
				if ok {
					t.Errorf("could find user in request context did not expect it")
				}
				return nil, nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := JWTAuthenticator{keyStorage: mockKeyStorage}
			if _, ok := tt.header["authorization"]; !ok {
				if tt.tokenGen.user != nil {
					token, err := gen.CreateToken(tt.tokenGen.user, tt.tokenGen.exp, tt.tokenGen.key)
					if err != nil {
						t.Fatalf("could not generate token: %v", err)
					}
					tt.header["authorization"] = fmt.Sprintf("Bearer %s", token)
				}
			}
			ctx := context.Background()
			ctx = metadata.NewIncomingContext(ctx, metadata.New(tt.header))

			_, err := InitInterceptor(gen)(ctx, new(interface{}), tt.info, mockHandler)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err %v;wanted %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("got unexpected err: %v", err)
			}
		})
	}
}
