package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/google/go-cmp/cmp"
	isPkg "github.com/matryer/is"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"testing"
	"time"
)

func TestUserModel_Create(t *testing.T) {
	isIntegrationTest(t)
	is := isPkg.New(t)
	wantName, wantEmail, wantPassword := "test", "test@example.org", "test123!"
	t.Run("inserts new user to database", func(t *testing.T) {
		is := is.New(t)
		db, teardown := getTestDb(t)
		defer teardown()

		model := UserModel{
			db: db,
		}

		err := model.Create(wantName, wantEmail, wantPassword)
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

	t.Run("returns error if user with same email exists", func(t *testing.T) {
		db, teardown := getTestDb(t)
		defer teardown()

		model := UserModel{
			db: db,
		}

		// insert first user with same fields than insert again to test duplicate email errors
		_ = model.Create(wantName, wantEmail, wantPassword)
		err := model.Create(wantName, wantEmail, wantPassword)
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
	isIntegrationTest(t)
	t.Run("returns user struct if user with id exists", func(t *testing.T) {
		db, teardown := dbInitializedForUsers(t)
		defer teardown()

		want := &models.User{
			ID:      1,
			Name:    "test",
			Email:   "test@example.org",
			Created: time.Date(2003, 8, 14, 18, 0, 0, 0, time.UTC),
		}

		model := &UserModel{
			db: db,
		}

		got, err := model.Get(1)
		if err != nil {
			t.Errorf("got error from getting in usermodel, did not expect one %v", err)
		}

		if !cmp.Equal(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("returns ErrNotFound if no user with id is found", func(t *testing.T) {
		db, teardown := getTestDb(t)
		defer teardown()

		model := &UserModel{
			db: db,
		}

		got, err := model.Get(1)
		if got != nil {
			t.Errorf("got user struct, did not expect one %v", got)
		}

		if err != models.ErrNotFound {
			t.Errorf("wrong error got %v but wanted %v", err, models.ErrNotFound)
		}
	})
}

func dbInitializedForUsers(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	db, teardown := getTestDb(t)
	setupScript, _ := ioutil.ReadFile("../testdata/user.sql")
	_, err := db.Exec(string(setupScript))
	if err != nil {
		t.Fatalf("got error initializing user into database: %v", err)
	}
	return db, teardown
}

func TestUserModel_Authenticate(t *testing.T) {
	isIntegrationTest(t)
	db, teardown := dbInitializedForUsers(t)
	defer teardown()
	model := &UserModel{
		db: db,
	}
	tests := []struct {
		email     string
		password  string
		wantedId  int
		wantErr   bool
		wantedErr error
	}{
		{
			email:     "test@example.org",
			password:  "password123",
			wantedId:  1,
			wantErr:   false,
			wantedErr: nil,
		},
		{
			email:     "test2@example.org",
			password:  "password123",
			wantedId:  2,
			wantErr:   false,
			wantedErr: nil,
		},
		{
			email:     "test1@example.org",
			password:  "password123",
			wantedId:  0,
			wantErr:   true,
			wantedErr: models.ErrInvalidCredentials,
		},
		{
			email:     "test@example.org",
			password:  "password",
			wantedId:  0,
			wantErr:   true,
			wantedErr: models.ErrInvalidCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("authenticate user with %q and %q", tt.email, tt.password), func(t *testing.T) {
			userId, err := model.Authenticate(tt.email, tt.password)
			if tt.wantErr {
				if err != nil && err != tt.wantedErr {
					t.Errorf("got err: %v but wanted %v", err, tt.wantedErr)
				}
			}
			if userId != tt.wantedId {
				t.Errorf("got userId %d but wanted %d", userId, tt.wantedId)
			}
		})
	}
}
