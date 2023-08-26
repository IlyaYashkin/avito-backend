package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func Open() *sql.DB {
	pgdb, err := sql.Open("postgres", getPsqlInfo())
	if err != nil {
		panic(err)
	}

	db = pgdb

	return pgdb
}

func Get() *sql.DB {
	return db
}

func getPsqlInfo() string {
	dbMap := map[string]string{
		"host":     os.Getenv("DB_HOST"),
		"user":     os.Getenv("DB_USER"),
		"password": os.Getenv("DB_PASSWORD"),
		"dbname":   os.Getenv("DB_NAME"),
		"port":     os.Getenv("DB_PORT"),
		"sslmode":  os.Getenv("DB_SSLMODE"),
		"TimeZone": os.Getenv("DB_TIMEZONE"),
	}

	var psqlInfo string

	for key, element := range dbMap {
		psqlInfo += fmt.Sprintf("%s=%s ", key, element)
	}

	return psqlInfo
}
