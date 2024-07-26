package dbpg

func (db *Postgres) Index() {
	query := "CREATE INDEX idx_json_data_1 ON json_data((data->>'key'))"

	tx, err := db.conn.Begin()
	if err != nil {
		panic("Cannot start transaction: " + err.Error())
	}

	_, err = tx.Exec(query)
	if err != nil {
		panic("Cannot create index: " + err.Error())
	}
	_ = tx.Commit()
}
