package db

import (
	"fmt"

	"github.com/codeboyzhou/sql-copilot/strconst"
)

type ExplainResult struct {
	ID           int64
	SelectType   string
	Table        string
	Partitions   string
	Type         string
	PossibleKeys string
	Key          string
	KeyLen       int64
	Ref          string
	Rows         int64
	Filtered     string
	Extra        string
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
