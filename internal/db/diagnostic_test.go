package db

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/codeboyzhou/sql-copilot/strconst"
)

const mockExplainSQL = "EXPLAIN SELECT * FROM table_example_user WHERE id = 1"

var mockExplainResult = ExplainResult{
	ID:           1,
	SelectType:   "SIMPLE",
	Table:        "table_example_user",
	Partitions:   sql.NullString{String: strconst.Empty, Valid: false},
	Type:         "ALL",
	PossibleKeys: sql.NullString{String: strconst.Empty, Valid: false},
	Key:          sql.NullString{String: strconst.Empty, Valid: false},
	KeyLen:       sql.NullInt64{Int64: 0, Valid: false},
	Ref:          sql.NullString{String: strconst.Empty, Valid: false},
	Rows:         1,
	Filtered:     "100.00",
	Extra:        sql.NullString{String: strconst.Empty, Valid: false},
}

type diagnosticMockRow struct{}

func (row *diagnosticMockRow) Scan(dest ...any) error {
	// mock the scan result for each field
	if len(dest) == 12 {
		// assign values to each field pointer
		if id, ok := dest[0].(*int64); ok {
			*id = mockExplainResult.ID
		}
		if selectType, ok := dest[1].(*string); ok {
			*selectType = mockExplainResult.SelectType
		}
		if table, ok := dest[2].(*string); ok {
			*table = mockExplainResult.Table
		}
		if partitions, ok := dest[3].(*sql.NullString); ok {
			*partitions = mockExplainResult.Partitions
		}
		if type_, ok := dest[4].(*string); ok {
			*type_ = mockExplainResult.Type
		}
		if possibleKeys, ok := dest[5].(*sql.NullString); ok {
			*possibleKeys = mockExplainResult.PossibleKeys
		}
		if key, ok := dest[6].(*sql.NullString); ok {
			*key = mockExplainResult.Key
		}
		if keyLen, ok := dest[7].(*sql.NullInt64); ok {
			*keyLen = mockExplainResult.KeyLen
		}
		if ref, ok := dest[8].(*sql.NullString); ok {
			*ref = mockExplainResult.Ref
		}
		if rows, ok := dest[9].(*int64); ok {
			*rows = mockExplainResult.Rows
		}
		if filtered, ok := dest[10].(*string); ok {
			*filtered = mockExplainResult.Filtered
		}
		if extra, ok := dest[11].(*sql.NullString); ok {
			*extra = mockExplainResult.Extra
		}
	}
	return nil
}

type diagnosticMockQuerier struct{}

func (querier *diagnosticMockQuerier) QueryRow(query string, args ...any) Row {
	return &diagnosticMockRow{}
}

func TestExplainSQL(t *testing.T) {
	tests := []struct {
		name    string
		querier Querier
		sql     string
		args    []any
		want    *ExplainResult
		wantErr bool
	}{
		{
			name:    "invalid sql",
			querier: &diagnosticMockQuerier{},
			sql:     strconst.Empty,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "explain sql",
			querier: &diagnosticMockQuerier{},
			sql:     mockExplainSQL,
			want:    &mockExplainResult,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := ExplainSQL(tt.querier, tt.sql)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExplainSQL(%s) failed, we don't want error, but got error: %v", tt.sql, gotErr)
				}
				// nothing to verify if we want error
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExplainSQL(%s) failed, got = %v, but want = %v", tt.sql, got, tt.want)
			}
		})
	}
}
