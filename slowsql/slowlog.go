package slowsql

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/codeboyzhou/sql-copilot/slowsql/strconst"
)

const (
	TimePrefix         = "# Time:"
	UserHostPrefix     = "# User@Host:"
	QueryTimePrefix    = "# Query_time:"
	SetTimestampPrefix = "SET timestamp="
)

type SlowQuery struct {
	QueryTime    float64
	RowsExamined int64
	SQL          string
}

func ParseSlowLog(filepath string, threshold float64) ([]SlowQuery, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var slowQueries []SlowQuery
	var currentQuery SlowQuery
	var sqlLines []string
	inSQLBlock := false

	// Precompile regex patterns for efficiency
	queryTimePattern := regexp.MustCompile(`Query_time:\s+([0-9.]+)`)
	rowsExaminedPattern := regexp.MustCompile(`Rows_examined:\s+(\d+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip unnecessary header lines
		if strings.HasPrefix(line, TimePrefix) || strings.HasPrefix(line, UserHostPrefix) || strings.HasPrefix(line, SetTimestampPrefix) {
			continue
		}

		// Process query metadata lines
		if strings.HasPrefix(line, QueryTimePrefix) {
			extractQueryTimeAndSet(&currentQuery, queryTimePattern, line)
			extractRowsExaminedAndSet(&currentQuery, rowsExaminedPattern, line)
			inSQLBlock = false
			continue
		}

		// Skip other comment lines
		if strings.HasPrefix(line, strconst.HashSymbol) {
			inSQLBlock = false
			continue
		}

		// Process SQL lines
		if line != strconst.Empty {
			inSQLBlock = true
			sqlLines = append(sqlLines, line)
		} else if inSQLBlock {
			// End of SQL block
			finalizeQuery(&currentQuery, sqlLines, threshold, &slowQueries)
			currentQuery = SlowQuery{}
			sqlLines = nil
			inSQLBlock = false
		}
	}

	// Handle case where file ends without empty line
	if inSQLBlock {
		finalizeQuery(&currentQuery, sqlLines, threshold, &slowQueries)
	}

	return slowQueries, nil
}

func extractQueryTimeAndSet(query *SlowQuery, queryTimePattern *regexp.Regexp, line string) {
	if matches := queryTimePattern.FindStringSubmatch(line); len(matches) > 1 {
		if queryTime, err := strconv.ParseFloat(matches[1], 64); err == nil {
			query.QueryTime = queryTime
		}
	}
}

func extractRowsExaminedAndSet(query *SlowQuery, rowsExaminedPattern *regexp.Regexp, line string) {
	if matches := rowsExaminedPattern.FindStringSubmatch(line); len(matches) > 1 {
		if rowsExamined, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
			query.RowsExamined = rowsExamined
		}
	}
}

func finalizeQuery(query *SlowQuery, sqlLines []string, threshold float64, slowQueries *[]SlowQuery) {
	if query.QueryTime >= threshold {
		query.SQL = strings.Join(sqlLines, strconst.Space)
		*slowQueries = append(*slowQueries, *query)
	}
}
