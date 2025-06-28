package dbmysql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Items map[string]interface{}

func (a Items) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Items) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func (db *MySQL) Load(filename string) {
	db.loadJSON(db.rows, filename)

	if db.transform {
		db.loadSchema()
	}
}
