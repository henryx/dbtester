package dbmysql

func (db *MySQL) Index() {
query := "CREATE INDEX idx_test_1 ON test.test(json_key)"

	tx, err := db.conn.Begin()
	if err != nil {
		panic("Cannot start transaction: " + err.Error())
	}
	
	_, err = tx.Exec(query)
	if err != nil {
		panic("Cannot create index: " + err.Error())
	}
	tx.Commit()
}