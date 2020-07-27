package dbpg

func (db *Postgres) Find() int64 {
	var counter int64
	
	counter = 0
	//j := make(map[string]interface{})
	j := make([]byte, 0)

	query := `SELECT data FROM test WHERE data->>'key' = $1`

	row, err := db.conn.Query(query, "/books/OL17806216M")
	if err != nil {
		panic("Cannot execute query: " + err.Error())
	}

	for row.Next() {
		err = row.Scan(&j)
		if err != nil {
			panic("Cannot read value: " + err.Error())
		}
		counter++
	}
	return counter
}
