package slowsql

import (
	"bufio"
	"os"
	"strconv"
	"strings"
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
	var isSQLBlock bool
	var sqlLines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "# Time:") || strings.HasPrefix(line, "# User@Host:") {
			continue
		}

		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "# Query_time:") {
				parts := strings.Fields(line)
				partsLength := len(parts)
				for i, part := range parts {
					if part == "Query_time:" && i+1 < partsLength {
						if queryTime, err := strconv.ParseFloat(parts[i+1], 64); err == nil {
							currentQuery.QueryTime = queryTime
						}
					}
					if part == "Rows_examined:" && i+1 < partsLength {
						if rowsExamined, err := strconv.ParseInt(parts[i+1], 10, 64); err == nil {
							currentQuery.RowsExamined = rowsExamined
						}
					}
				}
			}
			isSQLBlock = false
		} else if line != "" && !strings.HasPrefix(line, "SET timestamp=") {
			isSQLBlock = true
			sqlLines = append(sqlLines, line)
		} else if isSQLBlock && (line == "" || strings.HasPrefix(line, "SET timestamp=")) {
			if currentQuery.QueryTime >= threshold {
				currentQuery.SQL = strings.Join(sqlLines, " ")
				slowQueries = append(slowQueries, currentQuery)
			}
			currentQuery = SlowQuery{}
			isSQLBlock = false
			sqlLines = nil
		}
	}

	// check if the SQL block is the last line of the file
	if isSQLBlock && currentQuery.QueryTime >= threshold {
		currentQuery.SQL = strings.Join(sqlLines, " ")
		slowQueries = append(slowQueries, currentQuery)
	}

	return slowQueries, nil
}
