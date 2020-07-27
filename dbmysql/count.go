package dbmysql

func (db *MySQL) Count() int64 {
	var count int64
	var err error

	query := "SELECT count(*) FROM test.test"
	row := db.conn.QueryRow(query)
	
	err = row.Scan(&count)
	if err != nil {
		panic("Cannot convert count result: " + err.Error())
	}

	return count
}