package server

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/auth"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/metadata"
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

		want := ErrNoRefreshToken
		if err != want {
			t.Errorf("got err:%v, expected %v", err, want)
		}

	})
}

type mockUserStorage struct {
	called       map[string]bool
	authenticate func(ctx context.Context, email, password string) (*api.User, error)
	getKey       func(ctx context.Context, userId int32) ([]byte, error)
	insertKey    func(ctx context.Context, userId int32, key []byte) error
	deleteKey    func(ctx context.Context, userId int32) error
}

func (m *mockUserStorage) Authenticate(ctx context.Context, email, password string) (*api.User, error) {
	m.called["auth"] = true
	return m.authenticate(ctx, email, password)
}

func (m *mockUserStorage) GetRefreshKey(ctx context.Context, userId int32) ([]byte, error) {
	m.called["get"] = true
	return m.getKey(ctx, userId)
}

func (m *mockUserStorage) InsertRefreshKey(ctx context.Context, userId int32, key []byte) error {
	if m.called["delete"] {
		return nil
	}
	m.called["insert"] = true
	return m.insertKey(ctx, userId, key)
}

func (m *mockUserStorage) DeleteRefreshKey(ctx context.Context, userId int32) error {

	m.called["delete"] = true
	return m.deleteKey(ctx, userId)
}

type mockGenerator struct {
	expTime func(d time.Duration) time.Time
	create  func(user *api.User, expirationTime time.Time, key []byte) (string, error)
	verify  func(token string, key []byte) (*api.User, error)
}

func (m *mockGenerator) ExpirationTime(duration time.Duration) time.Time {
	return m.expTime(duration)
}

func (m *mockGenerator) CreateToken(user *api.User, expirationTime time.Time, key []byte) (string, error) {
	return m.create(user, expirationTime, key)
}

func (m *mockGenerator) VerifyToken(token string, key []byte) (*api.User, error) {
	return m.verify(token, key)
}

func (m *mockGenerator) CreateRandomKey() []byte {
	return []byte("randomkey")
}

func TestUserServer_AuthenticateUser(t *testing.T) {
	testUser := &api.User{
		Id:    1,
		Name:  "testuser",
		Email: "test@user.com",
	}
	type returnErrs struct {
		storageAuth            error
		storageGetKey          error
		storageInsertKey       error
		storageDeleteKey       error
		generatorVerify        error
		generatorCreateAccess  error
		generatorCreateRefresh error
	}

	tests := []struct {
		name              string
		header            map[string]string
		want              *api.AuthenticateResponse
		wantErr           error
		storageReturnUser *api.User
		returnErr         *returnErrs
	}{
		{
			name:              "authenticate successfully",
			header:            map[string]string{"authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("user:passwd")))},
			storageReturnUser: testUser,
			want: &api.AuthenticateResponse{
				TokenType:    api.AuthenticateResponse_BEARER,
				AccessToken:  "thisIsAnAccessTokenForTests",
				RefreshToken: "thisIsARefreshTokenForTests",
				ExpiresIn:    time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC).Unix(),
			},
		},
		{
			name:    "if storage can not identify user, returns ErrUserOrPasswdWrong",
			header:  map[string]string{"authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("user:passwd")))},
			wantErr: ErrNameOrPasswdWrong,
			returnErr: &returnErrs{
				storageAuth: models.ErrInvalidCredentials,
			},
		},
		{
			name:    "no authorization header returns ErrNoAuthHeader",
			header:  map[string]string{},
			wantErr: auth.ErrNoAuthHeader,
		},
		{
			name:    "no basic authorization header returns ErrNoBasicAuth",
			header:  map[string]string{"authorization": "Bearer <token>"},
			wantErr: auth.ErrNoBasicAuth,
		},
		{
			name:    "no basic authorization header returns ErrNoBasicAuth",
			header:  map[string]string{"authorization": "Bearer <token>"},
			wantErr: auth.ErrNoBasicAuth,
		},
		{
			name:    "generator could not create access token returns ErrCouldNotAuthorize",
			header:  map[string]string{"authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("user:passwd")))},
			wantErr: auth.ErrCouldNotAuthorize,
			returnErr: &returnErrs{
				generatorCreateAccess: errors.New("verifyToken could not create token"),
			},
		},
		{
			name:    "generator could not create access token returns ErrCouldNotAuthorize",
			header:  map[string]string{"authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("user:passwd")))},
			wantErr: auth.ErrCouldNotAuthorize,
			returnErr: &returnErrs{
				generatorCreateAccess: errors.New("verifyToken could not create token"),
			},
		},
		{
			name:    "generator could not create refresh token returns ErrCouldNotAuthorize",
			header:  map[string]string{"authorization": fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte("user:passwd")))},
			wantErr: auth.ErrCouldNotAuthorize,
			returnErr: &returnErrs{
				generatorCreateRefresh: errors.New("verifyToken could not create token"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &userServer{
				storage: &mockUserStorage{
					called: make(map[string]bool),
					authenticate: func(ctx context.Context, email, password string) (user *api.User, err error) {
						if tt.returnErr != nil && tt.returnErr.storageAuth != nil {
							return nil, tt.returnErr.storageAuth
						}
						return tt.storageReturnUser, nil
					},
					insertKey: func(ctx context.Context, userId int32, key []byte) error {
						if testUser.Id != userId {
							t.Errorf("user id does not match %d != %d", testUser.Id, userId)
						}
						return nil
					},
				},
				tokenGenerator: &mockGenerator{
					expTime: func(d time.Duration) time.Time {
						return time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
					},
					create: func(user *api.User, expirationTime time.Time, key []byte) (s string, err error) {
						if string(key) != string(auth.AccessTokenKey) {
							if tt.returnErr != nil && tt.returnErr.generatorCreateRefresh != nil {
								return "", tt.returnErr.generatorCreateRefresh
							}
							return "thisIsARefreshTokenForTests", nil
						}
						if tt.returnErr != nil && tt.returnErr.generatorCreateAccess != nil {
							return "", tt.returnErr.generatorCreateAccess
						}
						return "thisIsAnAccessTokenForTests", nil
					},
				},
			}
			ctx := metadata.NewIncomingContext(context.Background(), metadata.New(tt.header))

			got, err := server.AuthenticateUser(ctx, &empty.Empty{})
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err: %v, expected: %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			assertAuthenticateReponse(t, got, tt.want)
		})
	}
}

func TestUserServer_createRefreshToken(t *testing.T) {
	createErr := errors.New("verifyToken could not create Token")
	unkownErr := errors.New("storage error")

	type returnErrs struct {
		createToken error
		insertKey   error
		deleteKey   error
	}
	tests := []struct {
		name    string
		want    string
		wantErr error
		errors  *returnErrs
	}{
		{
			name: "create refreshtoken",
			want: "thisIsARefreshTokenForTests",
		},
		{
			name:    "create token returns error",
			wantErr: createErr,
			errors: &returnErrs{
				createToken: createErr,
			},
		},
		{
			name:    "insert key returns unkown error",
			wantErr: unkownErr,
			errors: &returnErrs{
				insertKey: unkownErr,
			},
		},
		{
			name:    "delete key returns error",
			wantErr: unkownErr,
			errors: &returnErrs{
				insertKey: models.ErrUserHasRefreshKey,
				deleteKey: unkownErr,
			},
		},
		{
			name: "insert key returns ErrUserHasRefreshKey error",
			want: "thisIsARefreshTokenForTests",
			errors: &returnErrs{
				insertKey: models.ErrUserHasRefreshKey,
			},
		},
		{
			name: "insert key returns ErrRefreshKeyIsInUse error",
			want: "thisIsARefreshTokenForTests",
			errors: &returnErrs{
				insertKey: models.ErrRefreshKeyIsInUse,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			returnedInserErr := false
			server := &userServer{
				storage: &mockUserStorage{
					called: make(map[string]bool),
					insertKey: func(ctx context.Context, userId int32, key []byte) error {
						if userId != 1 {
							t.Errorf("wanted userId 1 got %d", userId)
						}
						if string(key) != "randomkey" {
							t.Errorf("got key %q, wanted %q", key, "randomkey")
						}
						if !returnedInserErr && tt.errors != nil && tt.errors.insertKey != nil {
							returnedInserErr = true
							return tt.errors.insertKey
						}
						return nil
					},
					deleteKey: func(ctx context.Context, userId int32) error {
						if userId != 1 {
							t.Errorf("wanted userId 1 got %d", userId)
						}
						if tt.errors != nil && tt.errors.deleteKey != nil {
							return tt.errors.deleteKey
						}
						return nil
					},
				},
				tokenGenerator: &mockGenerator{
					expTime: func(d time.Duration) time.Time {
						return time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
					},
					create: func(user *api.User, expirationTime time.Time, key []byte) (s string, err error) {
						if string(key) != "randomkey" {
							t.Errorf("got random key %q, wanted %q", string(key), "randomkey")
						}
						if tt.errors != nil && tt.errors.createToken != nil {
							return "", tt.errors.createToken
						}

						return "thisIsARefreshTokenForTests", nil
					},
				},
			}

			got, err := server.createRefreshToken(context.Background(), &api.User{
				Id:    1,
				Name:  "test",
				Email: "test@user.com",
			})

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err: %v, wanted %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("did not expect err: %v", err)
			}

			if got != tt.want {
				t.Errorf("got refreshToken %q, wanted %q", got, tt.want)
			}
		})
	}
}

func TestUserServer_LogoutUser(t *testing.T) {
	mockUser := &api.User{
		Id:    1,
		Name:  "test",
		Email: "test@example.com",
	}
	type returnErrs struct {
		deleteKey error
	}

	tests := []struct {
		name    string
		user    *api.User
		errors  *returnErrs
		wantErr error
	}{
		{
			name: "logout user",
			user: mockUser,
		},
		{
			name: "could not delete refresh key",
			user: mockUser,
			errors: &returnErrs{
				deleteKey: errors.New("storage delete error"),
			},
			wantErr: ErrCouldNotLogOut,
		},
		{
			name:    "without user in context",
			wantErr: auth.ErrCouldNotAuthorize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := &userServer{
				storage: &mockUserStorage{
					called: make(map[string]bool),
					deleteKey: func(ctx context.Context, userId int32) error {
						if tt.errors != nil && tt.errors.deleteKey != nil {
							return tt.errors.deleteKey
						}
						return nil
					},
				},
			}

			ctx := context.Background()
			ctx = context.WithValue(ctx, "user", tt.user)
			_, err := server.LogoutUser(ctx, &empty.Empty{})
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err: %v, wanted %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("did not expect err: %v", err)
			}

		})
	}
}

func TestUserServer_RefreshToken(t *testing.T) {

	mockUser := &api.User{
		Id:    1,
		Name:  "test",
		Email: "test@example.com",
	}
	type returnErrs struct {
		verify      error
		createToken error
		getKey      error
	}

	tests := []struct {
		name    string
		header  map[string]string
		user    *api.User
		want    *api.AuthenticateResponse
		wantErr error
		errors  *returnErrs
	}{
		{
			name: "refresh acces token",
			header: map[string]string{
				"x-refresh-token": "thisIsARefreshTokenForTests",
			},
			user: mockUser,
			want: &api.AuthenticateResponse{
				TokenType:    api.AuthenticateResponse_BEARER,
				AccessToken:  "thisIsAnAccessTokenForTests2",
				RefreshToken: "thisIsARefreshTokenForTests",
				ExpiresIn:    time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC).Unix(),
			},
		},
		{
			name:    "requires refresh token header",
			header:  map[string]string{},
			user:    mockUser,
			wantErr: ErrNoRefreshToken,
		},
		{
			name: "requires refresh key in db",
			header: map[string]string{
				"x-refresh-token": "thisIsARefreshTokenForTests",
			},
			user: mockUser,
			errors: &returnErrs{
				getKey: errors.New("could not access refresh key in storage"),
			},
			wantErr: auth.ErrCouldNotAuthorize,
		},
		{
			name: "requires valid refresh token",
			header: map[string]string{
				"x-refresh-token": "thisIsARefreshTokenForTests",
			},
			user: mockUser,
			errors: &returnErrs{
				verify: errors.New("could not verify refresh token"),
			},
			wantErr: auth.ErrCouldNotAuthorize,
		},
		{
			name: "returns error if token could not be created",
			header: map[string]string{
				"authorization":   "Bearer thisIsAnAccessTokenForTests",
				"x-refresh-token": "thisIsARefreshTokenForTests",
			},
			errors: &returnErrs{
				createToken: errors.New("could not create access token"),
			},
			wantErr: auth.ErrCouldNotAuthorize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := userServer{
				storage: &mockUserStorage{
					called: make(map[string]bool),
					getKey: func(ctx context.Context, userId int32) (bytes []byte, err error) {
						if tt.errors != nil && tt.errors.getKey != nil {
							return nil, tt.errors.getKey
						}
						return []byte("randomkey"), nil
					},
				},
				tokenGenerator: &mockGenerator{
					expTime: func(d time.Duration) time.Time {
						return time.Date(2000, 1, 2, 3, 4, 5, 0, time.UTC)
					},
					verify: func(token string, key []byte) (user *api.User, err error) {
						if tt.errors != nil && tt.errors.verify != nil {
							return nil, tt.errors.verify
						}

						return &api.User{
							Id:    1,
							Name:  "test",
							Email: "test@user.com",
						}, nil
					},
					create: func(user *api.User, expirationTime time.Time, key []byte) (s string, err error) {
						if tt.errors != nil && tt.errors.createToken != nil {
							return "", tt.errors.createToken
						}
						return "thisIsAnAccessTokenForTests2", nil
					},
				},
			}

			ctx := context.Background()
			ctx = context.WithValue(ctx, "user", tt.user)
			ctx = metadata.NewIncomingContext(ctx, metadata.New(tt.header))

			got, err := server.RefreshToken(ctx, &empty.Empty{})
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err: %v, expected: %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("did not expect error: %v", err)
			}

			assertAuthenticateReponse(t, got, tt.want)
		})
	}
}

func assertAuthenticateReponse(t *testing.T, got *api.AuthenticateResponse, want *api.AuthenticateResponse) {
	t.Helper()

	if got.AccessToken != want.AccessToken {
		t.Errorf("access token is not the expected: %s != %s", got.AccessToken, want.AccessToken)
	}

	if got.RefreshToken != want.RefreshToken {
		t.Errorf("refresh token is not the expected: %s != %s", got.RefreshToken, want.RefreshToken)
	}

	if got.ExpiresIn != want.ExpiresIn {
		t.Errorf("expiry date is not the expected: %d != %d", got.ExpiresIn, want.ExpiresIn)
	}

	if got.TokenType != want.TokenType {
		t.Errorf("token type is not the expected: %s != %s", got.TokenType, want.TokenType)
	}
}
