package auth

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

var generator = NewJwtGenerator()

func TestCreateRandomKey(t *testing.T) {
	t.Run("create some random keys", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			got := generator.CreateRandomKey()
			if len(got) != 32 {
				t.Errorf("createRandomKey should've returned 64 char long string, got %d", len(got))
			}
		}
	})
}

func TestExpirationTime(t *testing.T) {
	inFiveMinutes := time.Now().Add(5 * time.Minute).Round(time.Second)
	got := generator.ExpirationTime(5 * time.Minute).Round(time.Second)

	if !got.Equal(inFiveMinutes) {
		t.Errorf("got %v wanted %v", got, inFiveMinutes)
	}
}

func mockTimeStamp() *timestamp.Timestamp {
	t, _ := ptypes.TimestampProto(time.Date(2019, 1, 18, 17, 16, 15, 0, time.UTC))
	return t
}

func TestCreateToken(t *testing.T) {
	type args struct {
		user *api.User
		time time.Time
		key  []byte
	}
	tests := []struct {
		name  string
		input args
	}{
		{
			name: "create token",
			input: args{
				user: &api.User{
					Id:      1,
					Name:    "testuser1",
					Email:   "test@example.com",
					Created: mockTimeStamp(),
				},
				time: time.Now().Add(time.Minute),
				key:  []byte("teststring"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generator.CreateToken(tt.input.user, tt.input.time, tt.input.key)
			if err != nil {
				t.Fatalf("could not create token: %v", err)
			}
			tkn, err := jwt.Parse(got, func(token *jwt.Token) (i interface{}, err error) {
				return tt.input.key, nil
			})

			if err != nil {
				t.Fatalf("could not verify token: %v", err)
			}

			if !tkn.Valid {
				t.Errorf("token is not valid: %v", tkn.Claims.Valid())
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	mUser := &api.User{
		Id:      1,
		Name:    "testuser1",
		Email:   "test@example.com",
		Created: mockTimeStamp(),
	}
	type args struct {
		user *api.User
		time time.Time
		key  []byte
	}
	tests := []struct {
		name    string
		token   string
		input   args
		want    *api.User
		wantErr error
	}{
		{
			name: "valid token",
			input: args{
				user: mUser,
				time: time.Now().Add(time.Minute),
				key:  []byte("teststring"),
			},
			want: mUser,
		},
		{
			name: "expired token",
			input: args{
				user: mUser,
				time: time.Now().Add(-10 * time.Minute),
				key:  []byte("teststring"),
			},
			wantErr: func() error {
				err := new(jwt.ValidationError)
				err.Inner = fmt.Errorf("token is expired by 10m0s")
				err.Errors |= jwt.ValidationErrorExpired
				return err
			}(),
		},
		{
			name: "tempered token",
			input: args{
				key: []byte("teststring"),
			},
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoyLCJuYW1lIjoidGVzdHVzZXIyIiwiZW1haWwiOiJ0ZXN0MkBleGFtcGxlLmNvbSIsImNyZWF0ZWQiOnsic2Vjb25kcyI6MTU0NzgzMTc3NX19LCJleHAiOjMyNTAyNjk5ODg1LCJqdGkiOiI2ZDQ0YWU2ZjczMGIyYmFkMzBiZjcwYzc3NDU3NzZiYiIsInN1YiI6InVzZXJfdGVzdHVzZXIxXzEifQ.SdNDLw32Olxh2LPw0FZXhWaYufEE56jzKEKfiwzDWE8",
			wantErr: func() error {
				err := new(jwt.ValidationError)
				err.Inner = fmt.Errorf("signature is invalid")
				err.Errors |= jwt.ValidationErrorSignatureInvalid
				return err
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.token
			if tt.token == "" {
				var err error
				token, err = generator.CreateToken(tt.input.user, tt.input.time, tt.input.key)
				if err != nil {
					t.Fatalf("could not create token: %v", err)
				}
			}

			got, err := generator.VerifyToken(token, tt.input.key)
			if tt.wantErr != nil {
				if err, ok := err.(*jwt.ValidationError); ok {
					want := tt.wantErr.(*jwt.ValidationError)
					if err.Inner.Error() != want.Inner.Error() || err.Errors != want.Errors {
						t.Errorf("got err [%d]%v, wanted [%d]%v", err.Errors, err, want.Errors, want)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("got err: %v, did not expect one", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got user %v, expected %v", got, tt.want)
			}

		})
	}
}
