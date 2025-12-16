package db

import (
	"fmt"
	"regexp"
)

type ShowCreateTableResult struct {
	TableName      string
	CreateTableSQL string
}

func ShowCreateTableSQL(querier Querier, tableName string) (result *ShowCreateTableResult, err error) {
	// validate table name
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(tableName) {
		return nil, fmt.Errorf("invalid table name: %s", tableName)
	}

	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	row := querier.QueryRow(query)
	result = &ShowCreateTableResult{}
	if err := row.Scan(&result.TableName, &result.CreateTableSQL); err != nil {
		return nil, fmt.Errorf("error while scanning show create table result: %w", err)
	}
	return result, nil
}
