package mysql

import "database/sql"

// createNullableString returns sql.NullString from given string (str)
// if str is empty string ( == ""), sql.NullString will be null
func createNullableString(str string) sql.NullString {
	var nullString sql.NullString
	if str != "" {
		// error can't happen in a simple string scan
		_ = nullString.Scan(str)
	}
	return nullString
}

// decodeNullableString will return string value from nullString
// if nullString is not valid (null), decodeNullableString will return empty string (== "")
func decodeNullableString(nullString sql.NullString) string {
	str := ""

	if nullString.Valid {
		str = nullString.String
	}

	return str
}
