package dbaccess

import (
	"errors"
)

func InsertUser(userId int32, ex QueryExecutor) (bool, error) {
	result, err := ex.Exec("insert into users (id) values ($1) on conflict (id) do nothing", userId)
	if err != nil {
		return false, errors.New("insert user error")
	}
	numRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return numRows == 1, nil
}
