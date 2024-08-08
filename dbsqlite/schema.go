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

func nullString(v gjson.Result) sql.NullString {
	if v.Type == gjson.Null {
		return sql.NullString{}
	}

	return sql.NullString{
		String: v.String(),
		Valid:  true,
	}
}

func (db *SQLite) addEdition(tx *sql.Tx, j string) {
	var editionId int
	var err error
	var insert = "INSERT INTO editions(key, title, subtitle, format, publish_date, edition, description, pages) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id"

	key := gjson.Get(j, "key")
	if !key.Exists() {
		return
	}

	title := gjson.Get(j, "title")
	if !title.Exists() {
		return
	}

	subtitle := gjson.Get(j, "subtitle")
	format := gjson.Get(j, "format")
	publishDate := gjson.Get(j, "publish_date")
	edition := gjson.Get(j, "edition_name")

	desc := gjson.Get(j, "description")

	description := ""
	switch desc.Type {
	case gjson.String:
		description = desc.String()
	case gjson.JSON:
		description = desc.Get("description.value").String()
	}

	pages := gjson.Get(j, "number_of_pages")

	stmt, err := tx.Prepare(insert)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}
	defer stmt.Close()

	err = stmt.QueryRow(key.String(), title.String(), nullString(subtitle), nullString(format),
		nullString(publishDate), nullString(edition), description, nullString(pages)).Scan(&editionId)

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

func (db *SQLite) loadEditions(tx *sql.Tx) {
	var err error

	query := `SELECT data FROM json_data WHERE data->'type'->>'key' = ?`
	j := ""

	tx, err = db.conn.Begin()
	if err != nil {
		panic("Cannot open transaction: " + err.Error())
	}

	row, err := tx.Query(query, common.EDITIONS)
	if err != nil {
		panic("Cannot execute query: " + err.Error())
	}

	for row.Next() {
		err = row.Scan(&j)
		if err != nil {
			panic("Cannot read value: " + err.Error())
		}
		db.addEdition(tx, j)

		err = tx.Commit()
		if err != nil {
			panic("Cannot commit transaction: " + err.Error())
		}
	}
}

func (db *SQLite) loadSchema() {
	var tx *sql.Tx
	var err error

	tx, err = db.conn.Begin()
	if err != nil {
		panic("Cannot open transaction: " + err.Error())
	}

	db.loadAuthors(tx)

	err = tx.Commit()
	if err != nil {
		panic("Cannot commit transaction: " + err.Error())
	}

	tx, err = db.conn.Begin()
	if err != nil {
		panic("Cannot open transaction: " + err.Error())
	}

	db.loadEditions(tx)

	err = tx.Commit()
	if err != nil {
		panic("Cannot commit transaction: " + err.Error())
	}
}
