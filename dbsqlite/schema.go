package dbsqlite

import (
	"database/sql"
	"dbtest/common"
	"github.com/tidwall/gjson"
)

func nullInt(v gjson.Result) sql.NullInt64 {
	if v.Type == gjson.Null {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: v.Int(),
		Valid: true,
	}
}

func (db *SQLite) addAuthor(tx *sql.Tx, j string) {
	var err error
	var insert = "INSERT INTO authors(key, name, revision, birth, death) VALUES(?, ?, ?, ?, ?)"

	key := gjson.Get(j, "key")
	if !key.Exists() {
		return
	}

	name := gjson.Get(j, "name")
	if !name.Exists() {
		return
	}

	revision := gjson.Get(j, "revision")
	birth := gjson.Get(j, "birth_date")
	death := gjson.Get(j, "death_date")

	_, err = tx.Exec(insert, key.String(), name.String(), revision.Int(), nullInt(birth), nullInt(death))
	if err != nil {
		panic("Cannot execute query: " + err.Error())
	}
}

func (db *SQLite) loadAuthors(tx *sql.Tx) {
	var err error

	query := `SELECT data FROM json_data WHERE data->'type'->>'key' = ?`
	j := ""

	row, err := tx.Query(query, common.AUTHORS)
	if err != nil {
		panic("Cannot execute query: " + err.Error())
	}

	for row.Next() {
		err = row.Scan(&j)
		if err != nil {
			panic("Cannot read value: " + err.Error())
		}
		db.addAuthor(tx, j)
	}
}

func (db *SQLite) loadSchema() {
	tx, err := db.conn.Begin()
	if err != nil {
		panic("Cannot open transaction: " + err.Error())
	}

	db.loadAuthors(tx)

	err = tx.Commit()
	if err != nil {
		panic("Cannot commit transaction: " + err.Error())
	}
}
