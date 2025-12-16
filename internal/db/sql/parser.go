package sql

import (
	"fmt"
	"slices"
	"strings"
	"unicode"

	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/codeboyzhou/sql-copilot/strconst"
)

func NormalizeSQL(sql string) string {
	sql = strings.TrimSpace(sql)

	var normalized strings.Builder
	previousIsSpace := false
	for _, char := range sql {
		if unicode.IsSpace(char) {
			if !previousIsSpace {
				normalized.WriteRune(strconst.SpaceRune)
				previousIsSpace = true
			}
		} else {
			normalized.WriteRune(char)
			previousIsSpace = false
		}
	}
	return normalized.String()
}

func ExtractTableNames(sql string) (tableNames []string, err error) {
	parser, err := sqlparser.New(sqlparser.Options{})
	if err != nil {
		return nil, fmt.Errorf("error new sql parser: %w", err)
	}

	normalizedSQL := NormalizeSQL(sql)
	stmt, err := parser.Parse(normalizedSQL)
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
