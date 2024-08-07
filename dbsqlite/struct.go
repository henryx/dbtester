package dbsqlite

import (
	"database/sql"
	"dbtest/cli"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	conn     *sql.DB
	database string
}

func (db *SQLite) create() {
	query := "CREATE TABLE IF NOT EXISTS json_data(data JSONB)"

	tx, err := db.conn.Begin()
	if err != nil {
		panic("Cannot start transaction: " + err.Error())
	}

	_, err = tx.Exec(query)
	if err != nil {
		panic("Cannot create table: " + err.Error())
	}
	_ = tx.Commit()
}

func (db *SQLite) clean() {
	query := "DROP TABLE IF EXISTS json_data"

	tx, err := db.conn.Begin()
	if err != nil {
		panic("Cannot start transaction: " + err.Error())
	}

	_, err = tx.Exec(query)
	if err != nil {
		panic("Cannot drop table: " + err.Error())
	}
	_ = tx.Commit()
}

func (db *SQLite) New(cli *cli.CLI) {
	var err error

	db.database = cli.SQLite.Database

	dsn := fmt.Sprintf("file:%s?_journal=WAL&_fk=true", db.database)
	db.conn, err = sql.Open("sqlite3", dsn)

	if err != nil {
		panic("Cannot open database connection: " + err.Error())
	}

	db.clean()
	db.create()
}

func (db *SQLite) Close() {
	_ = db.conn.Close()
}

func (db *SQLite) Name() string {
	return "SQLite"
}

func (db *SQLite) Url() string {
	return db.database
}
