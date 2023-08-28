package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type QueryExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func Open() *sql.DB {
	pgdb, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	db = pgdb

	return pgdb
}

func Get() *sql.DB {
	return db
}
