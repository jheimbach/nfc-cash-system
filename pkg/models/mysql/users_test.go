package mysql

import (
	"database/sql"
	"errors"
	"github.com/JHeimbach/nfc-cash-system/pkg/models"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestUserModel_Insert(t *testing.T) {
	wantName, wantEmail, wantPassword := "test", "test@example.org", "test123!"

	t.Run("happy path check arguments", func(t *testing.T) {
		db, teardown := getTestDb(t)
		defer teardown()

		model := UserModel{
			db: db,
		}

		err := model.Insert(wantName, wantEmail, wantPassword)
		if err != nil {
			t.Fatalf("got error from inserting in usermodel, did not expect one %v", err)
		}
		gotName, gotEmail, gotPassword := "", "", ""

		err = db.QueryRow("SELECT name,email,hashed_password from users WHERE id=1").Scan(&gotName, &gotEmail, &gotPassword)
		if err != nil {
			t.Fatalf("got error from inserting in usermodel, did not expect one %v", err)
		}
		assertEqualStrings(t, gotName, wantName)
		assertEqualStrings(t, gotEmail, wantEmail)
		assertEqualPasswords(t, gotPassword, wantPassword)
	})

	t.Run("returns error, got duplicate email", func(t *testing.T) {
		db, teardown := getTestDb(t)
		defer teardown()

		model := UserModel{
			db: db,
		}

		model.Insert(wantName, wantEmail, wantPassword)
		err := model.Insert(wantName, wantEmail, wantPassword)
		if err == nil {
			t.Fatalf("got no error, expected one")
		}
		if !errors.Is(err, models.ErrDuplicateEmail) {
			t.Errorf("got error %v, expected %v", err, models.ErrDuplicateEmail)
		}
	})
}

func TestUserModel_Get(t *testing.T) {
	db, teardown := getTestDb(t)
	defer teardown()
	setupScript, _ := ioutil.ReadFile("../testdata/user_get.sql")
	db.Exec(string(setupScript))

	want := &models.User{
		ID:      1,
		Name:    "test",
		Email:   "test@example.org",
		Created: time.Date(2003, 8, 14, 18, 0, 0, 0, time.UTC),
	}

	model := UserModel{
		db: db,
	}

	got, err := model.Get(1)
	if err != nil {
		t.Errorf("got error from getting in usermodel, did not expect one %v", err)
	}

	if !cmp.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

}

func assertEqualStrings(t *testing.T, got string, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertEqualPasswords(t *testing.T, got, want string) {
	t.Helper()
	err := bcrypt.CompareHashAndPassword([]byte(got), []byte(want))
	if err != nil {
		t.Errorf("passwords dont match")
	}
}

func getTestDb(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skipf("no database dsn found, skipping test")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skipf("could not connect to database, skipping test, err: %v", err)
	}

	if err = db.Ping(); err != nil {
		t.Skipf("could not connect to database, skipping test, err: %v", err)
	}
	setup, _ := ioutil.ReadFile("../migrations/20191028204458_users.up.sql")
	teardown, _ := ioutil.ReadFile("../migrations/20191028204458_users.down.sql")
	db.Exec(string(setup))

	teardownF := func() {
		db.Exec(string(teardown))
		db.Close()
	}

	return db, teardownF
}
