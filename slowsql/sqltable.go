package slowsql

import (
	"regexp"
	"strings"

	"github.com/codeboyzhou/sql-copilot/strconst"
)

func ExtractTableNames(sql string) []string {
	sqlKeywordFromPattern := regexp.MustCompile(`(?i)FROM\s+([\w.]+)`)
	sqlKeywordJoinPattern := regexp.MustCompile(`(?i)JOIN\s+([\w.]+)`)

	fromMatches := sqlKeywordFromPattern.FindAllStringSubmatch(sql, -1)
	joinMatches := sqlKeywordJoinPattern.FindAllStringSubmatch(sql, -1)

	var tableNames []string

	for _, match := range fromMatches {
		if len(match) > 1 {
			tableName := match[1]
			if strings.Contains(tableName, strconst.DOT) {
				tableName = strings.Split(tableName, strconst.DOT)[1]
			}
			tableNames = append(tableNames, tableName)
		}
	}

	for _, match := range joinMatches {
		if len(match) > 1 {
			tableName := match[1]
			if strings.Contains(tableName, strconst.DOT) {
				tableName = strings.Split(tableName, strconst.DOT)[1]
			}
			tableNames = append(tableNames, tableName)
		}
	}

	return tableNames
}
