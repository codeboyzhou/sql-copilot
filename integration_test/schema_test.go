// integration_test/schema_test.go
//go:build integration
// +build integration

package integration

import (
	"os"
	"testing"

	"github.com/codeboyzhou/sql-copilot/internal/db"
	"github.com/codeboyzhou/sql-copilot/strconst"
)

const wantTableName = "sale"

func TestShowCreateTableSQL(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == strconst.Empty {
		t.Fatalf("TEST_DB_DSN environment variable is not set")
	}

	dbConnection := db.Connect(dsn)
	defer dbConnection.Close()

	querier := db.NewQuerier(dbConnection)
	result, err := db.ShowCreateTableSQL(querier, wantTableName)
	if err != nil {
		t.Fatalf("Failed to show create table SQL: %v", err)
	}

	if result.TableName != wantTableName {
		t.Errorf("TestShowCreateTableSQL(%s) failed, got table name = %s, but want %s", wantTableName, result.TableName, wantTableName)
	}

	if result.CreateTableSQL == strconst.Empty {
		t.Errorf("TestShowCreateTableSQL(%s) failed, got create table SQL = %s, but want non-empty string", wantTableName, result.CreateTableSQL)
	}
}
