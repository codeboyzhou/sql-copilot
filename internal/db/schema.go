package db

import (
	"fmt"
	"regexp"

	"github.com/codeboyzhou/sql-copilot/strconst"
)

func ShowCreateTableSQL(querier Querier, tableName string) (createTableSQL string, err error) {
	// validate table name
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(tableName) {
		return strconst.Empty, fmt.Errorf("invalid table name: %s", tableName)
	}

	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
	row := querier.QueryRow(query)
	if err := row.Scan(&createTableSQL); err != nil {
		return strconst.Empty, fmt.Errorf("error while scanning show create table result: %w", err)
	}
	return createTableSQL, nil
}
