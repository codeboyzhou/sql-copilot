package sqlmeta

import (
	"fmt"
	"regexp"
	"slices"

	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/codeboyzhou/sql-copilot/internal/db"
	"github.com/codeboyzhou/sql-copilot/strconst"
)

func ExtractTableNames(sql string) (tableNames []string, err error) {
	parser, err := sqlparser.New(sqlparser.Options{})
	if err != nil {
		return nil, fmt.Errorf("error new sql parser: %w", err)
	}

	stmt, err := parser.Parse(sql)
	if err != nil {
		return nil, fmt.Errorf("error while parsing sql [%s]: %w", sql, err)
	}

	walkErr := sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch n := node.(type) {
		case sqlparser.TableName:
			if name := n.Name.String(); name != strconst.Empty && !slices.Contains(tableNames, name) {
				tableNames = append(tableNames, name)
			}
		}
		return true, nil
	}, stmt)
	if walkErr != nil {
		return nil, fmt.Errorf("error while walking sql AST: %w", walkErr)
	}

	return tableNames, nil
}

func GetTableSchema(querier db.Querier, tableName string) (createTableSQL string, err error) {
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
