package slowlog

import (
	"reflect"
	"testing"
)

func TestParseSlowLog(t *testing.T) {
	tests := []struct {
		name      string
		filepath  string
		threshold float64
		want      []SlowQuery
		wantErr   bool
	}{
		{
			name:      "should parse success",
			filepath:  "testdata/slowlog_test.txt",
			threshold: 1.0,
			want: []SlowQuery{
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
			want:      []SlowQuery{},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := ParseSlowLog(tt.filepath, tt.threshold)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ParseSlowLog(%s) failed, we don't want error, but got error: %v", tt.filepath, gotErr)
				}
				// nothing to verify if we want error
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSlowLog(%s), got = %v, but want = %v", tt.filepath, got, tt.want)
			}
		})
	}
}
