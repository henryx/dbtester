package dbmysql

import "fmt"

func (db *MySQL) Count() int64 {
	var count int64
	var err error

	query := fmt.Sprintf("SELECT count(*) FROM %s.json_data", db.database)
	row := db.conn.QueryRow(query)

	err = row.Scan(&count)
	if err != nil {
		panic("Cannot convert count result: " + err.Error())
	}

	return count
}
