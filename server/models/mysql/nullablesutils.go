package mysql

import "database/sql"

func createNullableString(str string) sql.NullString {
	var nullString sql.NullString
	if str != "" {
		// error can't happen in a simple string scan
		_ = nullString.Scan(str)
	}
	return nullString
}

func decodeNullableString(nullString sql.NullString) string {
	str := ""

	if nullString.Valid {
		str = nullString.String
	}

	return str
}
