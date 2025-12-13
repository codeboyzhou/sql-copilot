package slowsql_test

import (
	"testing"

	"github.com/codeboyzhou/sql-copilot/slowsql"
)

func TestParseSlowLog(t *testing.T) {
	tests := []struct {
		name      string
		filepath  string
		threshold float64
		want      []slowsql.SlowQuery
		wantErr   bool
	}{
		{
			name:      "should parse successfully and get 2 slow query",
			filepath:  "testdata/slowlog_test.txt",
			threshold: 1.0,
			want: []slowsql.SlowQuery{
				{
					QueryTime:    2.345,
					RowsExamined: 500000,
					SQL:          "SELECT * FROM orders WHERE user_id = 123 AND created > '2025-01-01';",
				},
				{
					QueryTime:    1.234,
					RowsExamined: 100000,
					SQL:          "SELECT * FROM orders WHERE user_id = 456 AND created > '2025-10-01';",
				},
			},
			wantErr: false,
		},
		{
			name:      "should parse error",
			filepath:  "testdata/slowlog_test_fake.txt",
			threshold: 1.0,
			want:      []slowsql.SlowQuery{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := slowsql.ParseSlowLog(tt.filepath, tt.threshold)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseSlowLog(%s) failed, we don't want error, but got: %v", tt.filepath, gotErr)
				}
				// nothing to verify if we want error
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("ParseSlowLog(%s), len(got) = %d, but len(want) = %d", tt.filepath, len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i].QueryTime != tt.want[i].QueryTime {
					t.Errorf("ParseSlowLog(%s), got[%d].QueryTime = %f, but want[%d].QueryTime = %f",
						tt.filepath, i, got[i].QueryTime, i, tt.want[i].QueryTime)
				}
				if got[i].RowsExamined != tt.want[i].RowsExamined {
					t.Errorf("ParseSlowLog(%s), got[%d].RowsExamined = %d, but want[%d].RowsExamined = %d",
						tt.filepath, i, got[i].RowsExamined, i, tt.want[i].RowsExamined)
				}
				if got[i].SQL != tt.want[i].SQL {
					t.Errorf("ParseSlowLog(%s), got[%d].SQL = %s, but want[%d].SQL = %s",
						tt.filepath, i, got[i].SQL, i, tt.want[i].SQL)
				}
			}
		})
	}
}
