package db

import "database/sql"

type Row interface {
	Scan(dest ...any) error
}

type Querier interface {
	QueryRow(query string, args ...any) Row
}

type row struct {
	r *sql.Row
}

func (r *row) Scan(dest ...any) error {
	return r.r.Scan(dest...)
}

type querier struct {
	db *sql.DB
}

func NewQuerier(db *sql.DB) Querier {
	return &querier{db: db}
}

func (q *querier) QueryRow(query string, args ...any) Row {
	return &row{r: q.db.QueryRow(query, args...)}
}
