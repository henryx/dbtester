package dbpg

func (db *Postgres) Index() {
query := "CREATE INDEX idx_test_1 ON test((data->>'key'))"

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