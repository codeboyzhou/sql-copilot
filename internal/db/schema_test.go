package db

import (
	"reflect"
	"testing"

	"github.com/codeboyzhou/sql-copilot/internal/db/sql"
)

var mockShowCreateTableResult = &ShowCreateTableResult{
	TableName: "table_example_user",
	CreateTableSQL: `
        CREATE TABLE table_example_user (
		    id INT PRIMARY KEY,
		    name VARCHAR(255),
		    email VARCHAR(255),
		    age INT,
		    INDEX idx_email (email)
		)`,
}

type schemaMockRow struct{}

func (row *schemaMockRow) Scan(dest ...any) error {
	// mock the scan result for each field
	if len(dest) == 2 {
		// assign values to each field pointer
		if tableName, ok := dest[0].(*string); ok {
			*tableName = mockShowCreateTableResult.TableName
		}
		if createTableSQL, ok := dest[1].(*string); ok {
			*createTableSQL = mockShowCreateTableResult.CreateTableSQL
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
		want      *ShowCreateTableResult
		wantErr   bool
	}{
		{
			name:      "invalid table name",
			querier:   &schemaMockQuerier{},
			tableName: "invalid#table#name",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "table example user",
			querier:   &schemaMockQuerier{},
			tableName: "table_example_user",
			want:      mockShowCreateTableResult,
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

			got.CreateTableSQL = sql.NormalizeSQL(got.CreateTableSQL)
			tt.want.CreateTableSQL = sql.NormalizeSQL(tt.want.CreateTableSQL)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShowCreateTableSQL(%s) failed, got = %v, but want = %v", tt.tableName, got, tt.want)
			}
		})
	}
}
