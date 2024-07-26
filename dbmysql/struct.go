package dbmysql

import (
	"database/sql"
	"dbtest/cli"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	conn     *sql.DB
	host     string
	database string
}

func (db *MySQL) createTable() {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.json_data(data JSON, json_key VARCHAR(50) GENERATED ALWAYS AS (data->>'$.key'))", db.database)

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

func (db *MySQL) createDB() {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", db.database)

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

func (db *MySQL) clean() {
	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", db.database)

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

func (db *MySQL) New(cli *cli.CLI) {
	var err error

	db.host = cli.MySQL.Host
	port := cli.MySQL.Port
	user := cli.MySQL.User
	password := cli.MySQL.Password
	db.database = cli.MySQL.Database

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", user, password, db.host, port)
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
	return "MySQL"
}

func (db *MySQL) Url() string {
	return db.host
}
