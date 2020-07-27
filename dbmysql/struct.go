package dbmysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	conn *sql.DB
}

func (db *MySQL) createTable() {
	query := "CREATE TABLE IF NOT EXISTS test.test(data JSON, json_key VARCHAR(50) GENERATED ALWAYS AS (data->>'$.key'))"

	tx, err := db.conn.Begin()
	if err != nil {
		panic("Cannot start transaction: " + err.Error())
	}

	_, err = tx.Exec(query)
	if err != nil {
		panic("Cannot create table: " + err.Error())
	}
	tx.Commit()
}

func (db *MySQL) createDB() {
	query := "CREATE DATABASE IF NOT EXISTS test"

	tx, err := db.conn.Begin()
	if err != nil {
		panic("Cannot start transaction: " + err.Error())
	}

	_, err = tx.Exec(query)
	if err != nil {
		panic("Cannot create table: " + err.Error())
	}
	tx.Commit()
}

func (db *MySQL) clean() {
	query := "DROP DATABASE IF EXISTS test"

	tx, err := db.conn.Begin()
	if err != nil {
		panic("Cannot start transaction: " + err.Error())
	}

	_, err = tx.Exec(query)
	if err != nil {
		panic("Cannot drop table: " + err.Error())
	}
	tx.Commit()
}

func (db *MySQL) New(host string) {
	var err error

	dsn := fmt.Sprintf("root:Latina,1@tcp(%s:3306)/mysql", host)
	db.conn, err = sql.Open("mysql", dsn)
	if err != nil {
		panic("Cannot open database connection: " + err.Error())
	}

	db.clean()
	db.createDB()
	db.createTable()
}

func (db *MySQL) Close() {
	db.conn.Close()
}

func (db *MySQL) Name() string {
	return "MySQL/MariaDB"
}
