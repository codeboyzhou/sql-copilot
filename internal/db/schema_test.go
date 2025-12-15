package db

import (
	"testing"

	"github.com/codeboyzhou/sql-copilot/strconst"
)

const mockCreateTableSQL = `
    CREATE TABLE table_example_user (
		id INT PRIMARY KEY,
		name VARCHAR(255),
		email VARCHAR(255),
		age INT,
		INDEX idx_email (email),
	)`

type schemaMockRow struct{}

func (row *schemaMockRow) Scan(dest ...any) error {
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

type schemaMockQuerier struct{}

func (querier *schemaMockQuerier) QueryRow(query string, args ...any) Row {
	return &schemaMockRow{}
}

func TestShowCreateTableSQL(t *testing.T) {
	tests := []struct {
		name      string
		querier   Querier
		tableName string
		want      string
		wantErr   bool
	}{
		{
			name:      "invalid table name",
			querier:   &schemaMockQuerier{},
			tableName: "invalid#table#name",
			want:      strconst.Empty,
			wantErr:   true,
		},
		{
			name:      "table example user",
			querier:   &schemaMockQuerier{},
			tableName: "table_example_user",
			want:      mockCreateTableSQL,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := ShowCreateTableSQL(tt.querier, tt.tableName)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ShowCreateTableSQL(%s) failed, we don't want error, but got error: %v", tt.tableName, gotErr)
				}
				// nothing to verify if we want error
				return
			}

			if got != tt.want {
				t.Errorf("ShowCreateTableSQL(%s), got = %s, but want = %s", tt.tableName, got, tt.want)
			}
		})
	}
}
