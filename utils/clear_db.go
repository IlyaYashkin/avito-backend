package utils

import (
	"avito-backend/internal/database"
	"os"
)

// Call database.Open() before
func ClearDB() {
	c, err := os.ReadFile("/app/init.sql")
	if err != nil {
		panic(err)
	}
	sqlStr := string(c)

	db := database.Get()
	_, err = db.Exec(sqlStr)
	if err != nil {
		panic(err)
	}
}
