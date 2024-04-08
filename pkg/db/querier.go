package db

import "database/sql"

type Querier interface {
	Query(cmd string, args ...any) (*sql.Rows, error)
	QueryRow(cmd string, args ...any) *sql.Row
	Exec(cmd string, args ...any) (sql.Result, error)
}
