package dbaccess

import "errors"

func InsertUser(userId int32, ex QueryExecutor) error {
	_, err := ex.Exec("insert into users (id) values ($1) on conflict (id) do nothing", userId)
	if err != nil {
		return errors.New("insert user error")
	}
	return nil
}
