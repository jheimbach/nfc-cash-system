package mysql

import "database/sql"

// pageOffset returns offset from 0 based on page and size.
// if page is 0 or negative, pageOffset will return 0
func pageOffset(page int, size int) int {
	if page <= 0 {
		return 0
	}

	offset := (page - 1) * size
	return offset
}

// countAllIds returns int for given query.
// query should have style of SELECT COUNT(...) FROM ...
func countAllIds(db *sql.DB, query string, args ...interface{}) (int, error) {
	var count int
	err := db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// maxPageCount returns max number of pages that will contain all elements
// based on element count and size of page
func maxPageCount(count, size int) int {
	return (count + size - 1) / size
}
