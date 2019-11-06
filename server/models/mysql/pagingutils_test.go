package mysql

import (
	"database/sql"
	"testing"
)

func Test_countAllIds(t *testing.T) {
	isIntegrationTest(t)

	db, teardown := dbInitializedForAccountLists(t)
	defer teardown()

	type args struct {
		db    *sql.DB
		query string
		args  []interface{}
	}
	tests := []struct {
		args    args
		want    int
		wantErr bool
	}{
		{
			args: args{
				db:    db,
				query: "SELECT COUNT(id) FROM accounts ORDER BY id",
				args:  []interface{}{},
			},
			want: 9,
		},
		{
			args: args{
				db:    db,
				query: "SELECT COUNT(id) FROM accounts WHERE group_id=?",
				args:  []interface{}{1},
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		got, err := countAllIds(tt.args.db, tt.args.query, tt.args.args...)
		if (err != nil) != tt.wantErr {
			t.Errorf("countAllIds() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if got != tt.want {
			t.Errorf("countAllIds() got = %v, want %v", got, tt.want)
		}
	}
}

func Test_maxPageCount(t *testing.T) {
	tests := []struct {
		count int
		size  int
		want  int
	}{
		{
			count: 20,
			size:  10,
			want:  2,
		},
		{
			count: 19,
			size:  10,
			want:  2,
		},
		{
			count: 9,
			size:  10,
			want:  1,
		},
		{
			count: 9,
			size:  100,
			want:  1,
		},
	}
	for _, tt := range tests {
		if got := maxPageCount(tt.count, tt.size); got != tt.want {
			t.Errorf("maxPageCount(%d, %d) = %d, want %d", tt.count, tt.size, got, tt.want)
		}
	}
}

func Test_pageOffset(t *testing.T) {
	tests := []struct {
		page int
		size int
		want int
	}{
		{
			page: 0,
			size: 5,
			want: 0,
		},
		{
			page: 1,
			size: 5,
			want: 0,
		},
		{
			page: 1,
			size: 500,
			want: 0,
		},
		{
			page: 2,
			size: 5,
			want: 5,
		},
		{
			page: 2,
			size: 500,
			want: 500,
		},
	}
	for _, tt := range tests {
		if got := pageOffset(tt.page, tt.size); got != tt.want {
			t.Errorf("pageOffset(%d,%d) = %d, want %d", tt.page, tt.size, got, tt.want)
		}
	}
}
