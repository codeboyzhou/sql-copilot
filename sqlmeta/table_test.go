package sqlmeta

import (
	"reflect"
	"testing"
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
