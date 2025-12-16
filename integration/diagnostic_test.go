//go:build integration
// +build integration

package integration

import (
	"os"
	"testing"

	"github.com/codeboyzhou/sql-copilot/internal/db"
	"github.com/codeboyzhou/sql-copilot/strconst"
)

const sql = "SELECT * FROM sale"

func TestExplainSQL(t *testing.T) {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == strconst.Empty {
		t.Fatalf("TEST_DB_DSN environment variable is not set")
	}

	dbConnection := db.Connect(dsn)
	defer dbConnection.Close()

	querier := db.NewQuerier(dbConnection)
	result, err := db.ExplainSQL(querier, sql)
	if err != nil {
		t.Fatalf("Failed to explain SQL (%s): %v", sql, err)
	}

	if result.SelectType != "SIMPLE" {
		t.Errorf("TestExplainSQL(%s) failed, got select type = %s, but want = SIMPLE", sql, result.SelectType)
	}
	if result.Table != "sale" {
		t.Errorf("TestExplainSQL(%s) failed, got table = %s, but want = sale", sql, result.Table)
	}
	if result.Type != "ALL" {
		t.Errorf("TestExplainSQL(%s) failed, got type = %s, but want = ALL", sql, result.Type)
	}
}
