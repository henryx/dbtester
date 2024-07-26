package dbmysql

import "fmt"

func (db *MySQL) Index() {
	query := fmt.Sprintf("CREATE INDEX idx_json_data_1 ON %s.json_data(json_key)", db.database)

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
