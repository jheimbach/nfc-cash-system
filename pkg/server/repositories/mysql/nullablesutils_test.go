package mysql

import (
	"database/sql"
	"testing"
)

func Test_createNullableString(t *testing.T) {
	t.Run("create nullable string with empty string", func(t *testing.T) {
		got := createNullableString("")

		if got.Valid {
			t.Errorf("got valid string, expected null")
		}

		if got.String != "" {
			t.Errorf("got string value %q, expected empty string", got.String)
		}
	})

	t.Run("create nullable string with nonempty string", func(t *testing.T) {
		got := createNullableString("test")

		if !got.Valid {
			t.Errorf("got nil string, expected value")
		}

		if got.String != "test" {
			t.Errorf("got string value %q, expected %q", got.String, "test")
		}
	})
}

func Test_decodeNullableString(t *testing.T) {
	t.Run("decode null string", func(t *testing.T) {
		var nullStr sql.NullString
		got := decodeNullableString(nullStr)

		if got != "" {
			t.Errorf("got string value %q, expected empty string", got)
		}
	})

	t.Run("decode not null string", func(t *testing.T) {
		var nullStr sql.NullString
		_ = nullStr.Scan("test")
		got := decodeNullableString(nullStr)

		if got == "" {
			t.Errorf("got string value %q, expected empty string", got)
		}
	})
}
