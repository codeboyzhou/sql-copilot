package sqlmeta

import (
	"reflect"
	"testing"

	"github.com/codeboyzhou/sql-copilot/internal/db"
	"github.com/codeboyzhou/sql-copilot/strconst"
)

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
				t.Errorf("ExtractTableNames(%s), got = %v, but want = %v", tt.sql, got, tt.want)
			}
		})
	}
}

const mockCreateTableSQL = `
    CREATE TABLE table_example_user (
		id INT PRIMARY KEY,
		name VARCHAR(255),
		email VARCHAR(255),
		age INT,
		INDEX idx_email (email),
	)`

type mockRow struct{}

func (row *mockRow) Scan(dest ...any) error {
	// mock the scan result
	if len(dest) > 0 {
		// dest[0] should be a pointer to string
		switch value := dest[0].(type) {
		case *string:
			*value = mockCreateTableSQL
		}
	}
	return nil
}

type mockQuerier struct{}

func (querier *mockQuerier) QueryRow(query string, args ...any) db.Row {
	return &mockRow{}
}

func TestGetTableSchema(t *testing.T) {
	tests := []struct {
		name      string
		querier   db.Querier
		tableName string
		want      string
		wantErr   bool
	}{
		{
			name:      "invalid table name",
			querier:   &mockQuerier{},
			tableName: "invalid#table#name",
			want:      strconst.Empty,
			wantErr:   true,
		},
		{
			name:      "table example user",
			querier:   &mockQuerier{},
			tableName: "table_example_user",
			want:      mockCreateTableSQL,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := GetTableSchema(tt.querier, tt.tableName)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetTableSchema(%s) failed, we don't want error, but got error: %v", tt.tableName, gotErr)
				}
				// nothing to verify if we want error
				return
			}

			if got != tt.want {
				t.Errorf("GetTableSchema(%s), got = %s, but want = %s", tt.tableName, got, tt.want)
			}
		})
	}
}
