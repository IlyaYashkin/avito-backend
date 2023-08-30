package user

import (
	"avito-backend/internal/database"
	"errors"
)

type UserRepository interface {
	Save(id int32) (bool, error)
}

type UserRepositoryDB struct {
	ex database.QueryExecutor
}

func NewUserRepository(ex database.QueryExecutor) *UserRepositoryDB {
	return &UserRepositoryDB{ex: ex}
}

func (repo *UserRepositoryDB) Save(userId int32) (bool, error) {
	result, err := repo.ex.Exec("insert into users (id) values ($1) on conflict (id) do nothing", userId)
	if err != nil {
		return false, errors.New("insert user error")
	}
	numRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return numRows == 1, nil
}
