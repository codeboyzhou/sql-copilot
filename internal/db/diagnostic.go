package db

import (
	"database/sql"
	"fmt"

	"github.com/codeboyzhou/sql-copilot/strconst"
)

type ExplainResult struct {
	ID           int64
	SelectType   string
	Table        string
	Partitions   sql.NullString
	Type         string
	PossibleKeys sql.NullString
	Key          sql.NullString
	KeyLen       sql.NullInt64
	Ref          sql.NullString
	Rows         int64
	Filtered     string
	Extra        sql.NullString
}

func ExplainSQL(querier Querier, sql string, args ...any) (result *ExplainResult, err error) {
	// validate sql
	if sql == strconst.Empty {
		return nil, fmt.Errorf("sql is empty")
	}

	explainStatement := fmt.Sprintf("EXPLAIN %s", sql)
	row := querier.QueryRow(explainStatement, args...)
	if result, err = scanExplainResult(row); err != nil {
		return nil, err
	}
	return result, nil
}

func scanExplainResult(row Row) (result *ExplainResult, err error) {
	result = &ExplainResult{}
	if err := row.Scan(
		&result.ID,
		&result.SelectType,
		&result.Table,
		&result.Partitions,
		&result.Type,
		&result.PossibleKeys,
		&result.Key,
		&result.KeyLen,
		&result.Ref,
		&result.Rows,
		&result.Filtered,
		&result.Extra); err != nil {
		return nil, fmt.Errorf("error while scanning explain result: %w", err)
	}
	return result, nil
}
