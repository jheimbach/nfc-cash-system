package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/JHeimbach/nfc-cash-system/server/api"
	"github.com/JHeimbach/nfc-cash-system/server/models"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/go-cmp/cmp"
	isPkg "github.com/matryer/is"
	"golang.org/x/crypto/bcrypt"
)

func TestUserModel_Create(t *testing.T) {
	isIntegrationTest(t)
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
	isIntegrationTest(t)
	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	t.Run("returns userId struct if userId with id exists", func(t *testing.T) {
		dbSetup("../testdata/user.sql")
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
	isIntegrationTest(t)
	is := isPkg.New(t)
	db, dbSetup, dbTeardown := getTestDb(t)
	dbSetup("../testdata/user.sql")
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

func TestUserModel_InsertRefreshKey(t *testing.T) {
	isIntegrationTest(t)
	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	type args struct {
		userId     int32
		refreshKey []byte
	}
	tests := []struct {
		name         string
		wantErr      error
		insertBefore *args
		input        *args
	}{
		{
			name: "insert key",
			input: &args{
				userId:     1,
				refreshKey: []byte("55812817ad1f1baa775955ba2149443a551091c0561afe95c3d4ea796fcf38ec"),
			},
		},
		{
			name: "insert for userId id 1 a second key",
			input: &args{
				userId:     1,
				refreshKey: []byte("31722b526ec223ca98f1235613fc822117b551ef49c8b94ffd82e848cae25e6c"),
			},
			insertBefore: &args{
				userId:     1,
				refreshKey: []byte("55812817ad1f1baa775955ba2149443a551091c0561afe95c3d4ea796fcf38ec"),
			},
			wantErr: models.ErrUserHasRefreshKey,
		},
		{
			name: "insert same key to different userId",
			input: &args{
				userId:     2,
				refreshKey: []byte("55812817ad1f1baa775955ba2149443a551091c0561afe95c3d4ea796fcf38ec"),
			},
			insertBefore: &args{
				userId:     1,
				refreshKey: []byte("55812817ad1f1baa775955ba2149443a551091c0561afe95c3d4ea796fcf38ec"),
			},
			wantErr: models.ErrRefreshKeyIsInUse,
		},
		{
			name: "insert key for userId that does not exist",
			input: &args{
				userId:     100,
				refreshKey: []byte("55812817ad1f1baa775955ba2149443a551091c0561afe95c3d4ea796fcf38ec"),
			},
			wantErr: models.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbSetup("../testdata/user.sql")
			defer dbTeardown()

			if tt.insertBefore != nil {
				_, err := db.Exec("INSERT INTO users_refreshkeys (user_id, refresh_key) VALUES (?,?)", tt.insertBefore.userId, tt.insertBefore.refreshKey)
				if err != nil {
					t.Fatalf("could not insert before test, got err %v", err)
				}
			}

			model := UserModel{
				db: db,
			}

			err := model.InsertRefreshKey(context.Background(), tt.input.userId, tt.input.refreshKey)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err %v, expected %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("got err %v, did not expect one", err)
			}

			var userId int32
			var key string

			err = db.QueryRow(
				"SELECT user_id,refresh_key FROM users_refreshkeys WHERE user_id =?", tt.input.userId,
			).Scan(&userId, &key)
			if err != nil {
				t.Fatalf("got err %v, did not expect one", err)
			}

			if userId != tt.input.userId {
				t.Errorf("userid does not match: got %d, want %d", userId, tt.input.userId)
			}

			if key != string(tt.input.refreshKey) {
				t.Errorf("refreshkey does not match: got %q, want %q", key, tt.input.refreshKey)
			}

		})
	}
}

func TestUserModel_DeleteRefreshKey(t *testing.T) {
	isIntegrationTest(t)
	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	type args struct {
		userId     int32
		refreshKey string
	}
	tests := []struct {
		name         string
		wantErr      error
		insertBefore *args
		input        int32
	}{
		{
			name:  "delete key",
			input: 1,
			insertBefore: &args{
				userId:     1,
				refreshKey: "55812817ad1f1baa775955ba2149443a551091c0561afe95c3d4ea796fcf38ec",
			},
		},
		{
			name:  "delete key that does not exist",
			input: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbSetup("../testdata/user.sql")
			defer dbTeardown()

			if tt.insertBefore != nil {
				_, err := db.Exec("INSERT INTO users_refreshkeys (user_id, refresh_key) VALUES (?,?)", tt.insertBefore.userId, tt.insertBefore.refreshKey)
				if err != nil {
					t.Fatalf("could not insert before test, got err %v", err)
				}
			}

			model := UserModel{
				db: db,
			}

			err := model.DeleteRefreshKey(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err %v, expected %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("got err %v, did not expect one", err)
			}

			var userId int32
			var key string

			err = db.QueryRow(
				"SELECT user_id,refresh_key FROM users_refreshkeys WHERE user_id =?", tt.input,
			).Scan(&userId, &key)

			if err != sql.ErrNoRows {
				t.Errorf("should not have found row, got %v,%v", userId, key)
			}
		})
	}
}

func TestUserModel_GetRefreshKey(t *testing.T) {
	isIntegrationTest(t)
	db, dbSetup, dbTeardown := getTestDb(t)
	defer db.Close()

	tests := []struct {
		name        string
		wantErr     error
		want        string
		input       int32
		inserBefore bool
	}{
		{
			name:        "get key",
			want:        "55812817ad1f1baa775955ba2149443a551091c0561afe95c3d4ea796fcf38ec",
			input:       1,
			inserBefore: true,
		},
		{
			name:    "get key for user that has none",
			want:    "55812817ad1f1baa775955ba2149443a551091c0561afe95c3d4ea796fcf38ec",
			input:   1,
			wantErr: models.ErrNotFound,
		},
		{
			name:    "get key for user that does not exist",
			want:    "55812817ad1f1baa775955ba2149443a551091c0561afe95c3d4ea796fcf38ec",
			input:   1,
			wantErr: models.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbSetup("../testdata/user.sql")
			defer dbTeardown()

			if tt.inserBefore {
				_, err := db.Exec("INSERT INTO users_refreshkeys (user_id, refresh_key) VALUES (?,?)", tt.input, tt.want)
				if err != nil {
					t.Fatalf("could not insert before test, got err %v", err)
				}
			}

			model := UserModel{
				db: db,
			}

			got, err := model.GetRefreshKey(context.Background(), tt.input)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("got err %v, expected %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("got err %v, did not expect one", err)
			}

			if string(got) != tt.want {
				t.Errorf("refreshkey does not match: got %q, want %q", got, tt.want)
			}

		})
	}
}
