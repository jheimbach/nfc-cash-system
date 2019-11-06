package mysql

import "database/sql"

func pageOffset(page int, size int) int {
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * size
	return offset
}

func countAllIds(db *sql.DB, query string, args ...interface{}) (int, error) {
	var count int
	err := db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func maxPageCount(count, size int) int {
	return (count + size - 1) / size
}
