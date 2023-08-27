package dbaccess

func IsSegmentExists(name string, ex QueryExecutor) (bool, error) {
	var rowExists bool
	err := ex.QueryRow("select exists(select true from segments where name=$1)", name).Scan(&rowExists)
	if err != nil {
		return rowExists, err
	}
	return rowExists, nil
}

func InsertSegment(name string, userPercentage float32, ex QueryExecutor) error {
	_, err := ex.Exec("insert into segments (name, user_percentage) values ($1, $2)", name, userPercentage)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSegment(name string, ex QueryExecutor) error {
	_, err := ex.Exec("delete from segments values where name = $1", name)
	if err != nil {
		return err
	}
	return nil
}
