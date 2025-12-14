package slowlog

import (
	"fmt"
	"os"
	"testing"
)

func BenchmarkParseSlowLog(b *testing.B) {
	filepath := "testdata/slowlog_test.txt"
	threshold := 1.0

	for b.Loop() {
		_, err := ParseSlowLog(filepath, threshold)
		if err != nil {
			b.Fatalf("BenchmarkParseSlowLog failed: %v", err)
		}
	}
}

func BenchmarkParseSlowLogLarge(b *testing.B) {
	largeFile := "testdata/large_slowlog_test.txt"

	// Only generate the large file if it doesn't exist
	if _, err := os.Stat(largeFile); os.IsNotExist(err) {
		generateLargeTestFile(largeFile)
	}

	threshold := 0.5 // Lower threshold to include more queries

	for b.Loop() {
		_, err := ParseSlowLog(largeFile, threshold)
		if err != nil {
			b.Fatalf("BenchmarkParseSlowLogLarge failed: %v", err)
		}
	}
}

func generateLargeTestFile(filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		panic(fmt.Sprintf("Failed to create test file: %v", err))
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			panic(fmt.Sprintf("Failed to close test file: %v", closeErr))
		}
	}()

	// Generate 1000 slow queries
	for i := range 1000 {
		queryTime := 0.5 + float64(i%10)/10    // Varying query times between 0.5 and 1.4
		rowsExamined := 10000 * int64(i%100+1) // Varying rows examined

		if _, err := file.WriteString("# Time: 2025-12-11T23:30:00.123456Z\n"); err != nil {
			panic(fmt.Sprintf("Failed to write to test file: %v", err))
		}
		if _, err := file.WriteString("# User@Host: app[app] @ localhost []\n"); err != nil {
			panic(fmt.Sprintf("Failed to write to test file: %v", err))
		}
		if _, err := file.WriteString("# Thread_id: 12345  Schema: mydb  QC_hit: No\n"); err != nil {
			panic(fmt.Sprintf("Failed to write to test file: %v", err))
		}
		if _, err := fmt.Fprintf(file, "# Query_time: %.3f  Lock_time: 0.001  Rows_sent: 1000  Rows_examined: %d\n", queryTime, rowsExamined); err != nil {
			panic(fmt.Sprintf("Failed to write to test file: %v", err))
		}
		if _, err := file.WriteString("SET timestamp=1702336200;\n"); err != nil {
			panic(fmt.Sprintf("Failed to write to test file: %v", err))
		}
		if _, err := fmt.Fprintf(file, "SELECT * FROM orders WHERE user_id = %d AND created > '2025-01-01';\n\n", i); err != nil {
			panic(fmt.Sprintf("Failed to write to test file: %v", err))
		}
	}
}
