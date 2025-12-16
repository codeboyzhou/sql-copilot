package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/codeboyzhou/sql-copilot/strconst"
)

func Connect(dsn string) *sql.DB {
	db, err := sql.Open(strconst.MySQLDriverName, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
	}

	return db
}
