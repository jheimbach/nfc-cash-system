package mysql

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/internals/test"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/go-cmp/cmp"
	isPkg "github.com/matryer/is"
	"golang.org/x/crypto/bcrypt"
)

func TestUserModel_Create(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	wantName, wantEmail, wantPassword := "test", "test@example.org", "test123!"
	t.Run("inserts new userId to database", func(t *testing.T) {
		is := is.New(t)

		dbSetup()
		defer dbTeardown()

		model := UserModel{
			db: db,
		}

		err := model.Create(context.Background(), wantName, wantEmail, wantPassword)
		if err != nil {
			t.Fatalf("got error from inserting in usermodel, did not expect one %v", err)
		}
		gotName, gotEmail, gotPassword := "", "", ""

		err = db.QueryRow("SELECT name,email,hashed_password from users WHERE id=1").Scan(&gotName, &gotEmail, &gotPassword)
		if err != nil {
			t.Fatalf("got error from inserting in usermodel, did not expect one %v", err)
		}

		is.Equal(gotName, wantName)   // name is not the same
		is.Equal(gotEmail, wantEmail) // email is not the same
		assertEqualPasswords(t, gotPassword, wantPassword)
	})

	t.Run("returns error if userId with same email exists", func(t *testing.T) {
		dbSetup()
		defer dbTeardown()

		model := UserModel{
			db: db,
		}

		// insert first userId with same fields than insert again to test duplicate email errors
		_ = model.Create(context.Background(), wantName, wantEmail, wantPassword)
		err := model.Create(context.Background(), wantName, wantEmail, wantPassword)
		if err == nil {
			t.Fatalf("got no error, expected one")
		}
		if !errors.Is(err, models.ErrDuplicateEmail) {
			t.Errorf("got error %v, expected %v", err, models.ErrDuplicateEmail)
		}
	})
}

func assertEqualPasswords(t *testing.T, got, want string) {
	t.Helper()
	err := bcrypt.CompareHashAndPassword([]byte(got), []byte(want))
	if err != nil {
		t.Errorf("passwords dont match")
	}
}

func TestUserModel_Get(t *testing.T) {
	test.IsIntegrationTest(t)
	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	t.Run("returns userId struct if userId with id exists", func(t *testing.T) {
		dbSetup(dataFor("user"))
		defer dbTeardown()

		created, _ := ptypes.TimestampProto(time.Date(2003, 8, 14, 18, 0, 0, 0, time.UTC))

		want := &api.User{
			Id:      1,
			Name:    "test",
			Email:   "test@example.org",
			Created: created,
		}

		model := &UserModel{
			db: db,
		}

		got, err := model.Get(context.Background(), 1)
		if err != nil {
			t.Errorf("got error from getting in usermodel, did not expect one %v", err)
		}

		if !cmp.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("returns ErrNotFound if no userId with id is found", func(t *testing.T) {
		dbSetup()
		defer dbTeardown()

		model := &UserModel{
			db: db,
		}

		got, err := model.Get(context.Background(), 1)
		if got != nil {
			t.Errorf("got userId struct, did not expect one %v", got)
		}

		if err != models.ErrNotFound {
			t.Errorf("wrong error got %v but wanted %v", err, models.ErrNotFound)
		}
	})
}

func TestUserModel_Authenticate(t *testing.T) {
	test.IsIntegrationTest(t)
	is := isPkg.New(t)
	db, dbSetup, dbTeardown := getTestDb(t)
	dbSetup(dataFor("user"))
	defer func() {
		dbTeardown()
		db.Close()
	}()

	model := &UserModel{
		db: db,
	}

	mockUserOne := &api.User{
		Id:    1,
		Name:  "test",
		Email: "test@example.org",
		Created: func() *timestamp.Timestamp {
			t, _ := ptypes.TimestampProto(time.Date(2003, 8, 14, 18, 0, 0, 0, time.UTC))
			return t
		}(),
	}
	mockUserTwo := &api.User{
		Id:    2,
		Name:  "test",
		Email: "test2@example.org",
		Created: func() *timestamp.Timestamp {
			t, _ := ptypes.TimestampProto(time.Date(2003, 8, 14, 18, 0, 0, 0, time.UTC))
			return t
		}(),
	}

	tests := []struct {
		email    string
		password string
		want     *api.User
		wantErr  error
	}{
		{
			email:    "test@example.org",
			password: "password123",
			want:     mockUserOne,
			wantErr:  nil,
		},
		{
			email:    "test2@example.org",
			password: "password123",
			want:     mockUserTwo,
			wantErr:  nil,
		},
		{
			email:    "test1@example.org",
			password: "password123",
			want:     nil,
			wantErr:  models.ErrInvalidCredentials,
		},
		{
			email:    "test@example.org",
			password: "password",
			want:     nil,
			wantErr:  models.ErrInvalidCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("authenticate userId with %q and %q", tt.email, tt.password), func(t *testing.T) {
			is := is.New(t)
			got, err := model.Authenticate(context.Background(), tt.email, tt.password)
			if tt.wantErr != nil {
				is.Equal(err, tt.wantErr)
				return
			}

			is.Equal(got, tt.want)
		})
	}
}
