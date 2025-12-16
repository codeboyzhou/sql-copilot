package sql

import (
	"reflect"
	"testing"

	"github.com/codeboyzhou/sql-copilot/strconst"
)

const wantNormalizedSQL = "SELECT * FROM table1 WHERE id = 1"

func TestNormalizeSQL(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want string
	}{
		{
			name: "empty sql",
			sql:  strconst.Empty,
			want: strconst.Empty,
		},
		{
			name: "normalized sql",
			sql:  wantNormalizedSQL,
			want: wantNormalizedSQL,
		},
		{
			name: "sql with multiple spaces and tab spaces at the beginning",
			sql:  " \t \t SELECT * FROM table1 WHERE id = 1",
			want: wantNormalizedSQL,
		},
		{
			name: "sql with multiple spaces and tab spaces at the end",
			sql:  "SELECT * FROM table1 WHERE id = 1 \t \t ",
			want: wantNormalizedSQL,
		},
		{
			name: "sql with multiple spaces and tab spaces at the beginning and end",
			sql:  " \t \t SELECT * FROM table1 WHERE id = 1 \t \t ",
			want: wantNormalizedSQL,
		},
		{
			name: "sql is multiple lines",
			sql: `
			    SELECT *
				FROM table1
				WHERE id = 1
			`,
			want: wantNormalizedSQL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeSQL(tt.sql)
			if got != tt.want {
				t.Errorf("NormalizeSQL(%s) failed, got = %v, but want = %v", tt.sql, got, tt.want)
			}
		})
	}
}

func TestExtractTableNames(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		want    []string
		wantErr bool
	}{
		{
			name:    "simple from",
			sql:     "SELECT * FROM table1",
			want:    []string{"table1"},
			wantErr: false,
		},
		{
			name:    "simple join",
			sql:     "SELECT * FROM table1 LEFT JOIN table2 ON table1.id = table2.id",
			want:    []string{"table1", "table2"},
			wantErr: false,
		},
		{
			name:    "multiple joins",
			sql:     "SELECT * FROM table1 LEFT JOIN table2 ON table1.id = table2.id LEFT JOIN table3 ON table2.id = table3.id",
			want:    []string{"table1", "table2", "table3"},
			wantErr: false,
		},
		{
			name:    "case insensitive",
			sql:     "select * from table1 left join table2 on table1.id = table2.id",
			want:    []string{"table1", "table2"},
			wantErr: false,
		},
		{
			name:    "table with database prefix",
			sql:     "SELECT * FROM db.table1",
			want:    []string{"table1"},
			wantErr: false,
		},
		{
			name:    "no any from and join",
			sql:     "SELECT 1",
			want:    []string{"dual"},
			wantErr: false,
		},
		{
			name:    "not a sql statement",
			sql:     "This is not a sql statement",
			want:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := ExtractTableNames(tt.sql)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExtractTableNames(%s) failed, we don't want error, but got error: %v", tt.sql, gotErr)
				}
				// nothing to verify if we want error
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractTableNames(%s) failed, got = %v, but want = %v", tt.sql, got, tt.want)
			}
		})
	}
}
