package service

import (
	"database/sql"
	"errors"
)

func createUserUsingTransaction(userId int32, tx *sql.Tx) error {
	_, err := tx.Exec("insert into users (id) values ($1) on conflict (id) do nothing", userId)
	if err != nil {
		return errors.New("insert user error")
	}
	return nil
}
