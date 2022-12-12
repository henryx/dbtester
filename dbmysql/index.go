package dbmysql

import "fmt"

func (db *MySQL) Index() {
	query := fmt.Sprintf("CREATE INDEX idx_test_1 ON %s.test(json_key)", db.database)

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
