package slowsql

import (
	"testing"
)

func TestExtractTableNames(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "simple from",
			sql:  "SELECT * FROM table1",
			want: []string{"table1"},
		},
		{
			name: "simple join",
			sql:  "SELECT * FROM table1 LEFT JOIN table2 ON table1.id = table2.id",
			want: []string{"table1", "table2"},
		},
		{
			name: "multiple joins",
			sql:  "SELECT * FROM table1 LEFT JOIN table2 ON table1.id = table2.id LEFT JOIN table3 ON table2.id = table3.id",
			want: []string{"table1", "table2", "table3"},
		},
		{
			name: "no matches",
			sql:  "SELECT 1",
			want: []string{},
		},
		{
			name: "case insensitive",
			sql:  "select * from table1 left join table2 on table1.id = table2.id",
			want: []string{"table1", "table2"},
		},
		{
			name: "table with database prefix",
			sql:  "SELECT * FROM db.table1",
			want: []string{"table1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractTableNames(tt.sql)

			if len(got) != len(tt.want) {
				t.Errorf("ExtractTableNames(%s), len(got) = %d, but len(want) = %d", tt.sql, len(got), len(tt.want))
				return
			}

			for i, value := range got {
				if value != tt.want[i] {
					t.Errorf("ExtractTableNames(%s), got[%d] = %s, but want[%d] = %s", tt.sql, i, value, i, tt.want[i])
				}
			}
		})
	}
}
