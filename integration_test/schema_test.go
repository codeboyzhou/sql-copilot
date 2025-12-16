// integration_test/schema_test.go
//go:build integration
// +build integration

package integration_test

import (
	"os"
	"testing"

	"github.com/codeboyzhou/sql-copilot/internal/db"
	"github.com/codeboyzhou/sql-copilot/strconst"
)

func TestShowCreateTableSQL(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == strconst.Empty {
		t.Fatalf("TEST_DB_DSN environment variable is not set")
	}

	dbConnection := db.Connect(dsn)
	defer dbConnection.Close()

	querier := db.NewQuerier(dbConnection)
	result, err := db.ShowCreateTableSQL(querier, "sale")
	if err != nil {
		t.Fatalf("Failed to show create table SQL: %v", err)
	}
	if result.TableName != "sale" {
		t.Errorf("Expected table name 'sale', but got '%s'", result.TableName)
	}
	if result.CreateTableSQL == strconst.Empty {
		t.Errorf("Expected non-empty create table SQL, but got empty string")
	}
}
