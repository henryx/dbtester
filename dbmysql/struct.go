package dbmysql

import (
	"database/sql"
	"fmt"
	"gopkg.in/ini.v1"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	conn     *sql.DB
	host     string
	database string
}

func (db *MySQL) createTable() {
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.test(data JSON, json_key VARCHAR(50) GENERATED ALWAYS AS (data->>'$.key'))", db.database)

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
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", db.database)

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
	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", db.database)

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

func (db *MySQL) New(cfg *ini.Section) {
	var err error

	db.host = cfg.Key("host").MustString("localhost")
	port := cfg.Key("port").MustInt(3306)
	user := cfg.Key("user").MustString("root")
	password := cfg.Key("password").MustString("mysql")
	db.database = cfg.Key("database").MustString("libraries")

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
