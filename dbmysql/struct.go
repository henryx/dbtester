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
	tables := []string{
		fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.json_data(data JSON, json_key VARCHAR(50) GENERATED ALWAYS AS (data->>'$.key'))", db.database),
		"CREATE TABLE IF NOT EXISTS authors (id INTEGER PRIMARY KEY AUTOINCREMENT, key VARCHAR(60) NOT NULL, name VARCHAR(60) NOT NULL, revision INTEGER, birth INTEGER,death INTEGER, CONSTRAINT uk_authors_key UNIQUE (key, name))",
		"CREATE TABLE IF NOT EXISTS editions (id INTEGER PRIMARY KEY AUTOINCREMENT, key VARCHAR(60) NOT NULL, title VARCHAR(60) NOT NULL, subtitle VARCHAR(120), format TEXT, publish_date TEXT, edition TEXT, description TEXT, pages INTEGER, CONSTRAINT uk_editions_key UNIQUE (key, title))",
		"CREATE TABLE IF NOT EXISTS genres ( id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR(60) NOT NULL, CONSTRAINT uk_genres_key UNIQUE (name))",
		"CREATE TABLE IF NOT EXISTS publishers ( id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR(60) NOT NULL, CONSTRAINT uk_publishers_key UNIQUE (name))",
		"CREATE TABLE IF NOT EXISTS editions_genres ( id INTEGER PRIMARY KEY AUTOINCREMENT, edition_id INTEGER, genre_id INTEGER, CONSTRAINT fk_editions_genres_key_1 FOREIGN KEY (edition_id) REFERENCES editions (id), CONSTRAINT fk_editions_genres_key_2 FOREIGN KEY (genre_id) REFERENCES genres (id))",
		"CREATE TABLE IF NOT EXISTS editions_publishers ( id INTEGER PRIMARY KEY AUTOINCREMENT, edition_id INTEGER, publisher_id INTEGER, CONSTRAINT fk_editions_publishers_key_1 FOREIGN KEY (edition_id) REFERENCES editions (id), CONSTRAINT fk_editions_publishers_key_2 FOREIGN KEY (publisher_id) REFERENCES publishers (id))",
		"CREATE TABLE IF NOT EXISTS editions_isbn10 ( id INTEGER PRIMARY KEY AUTOINCREMENT, isbn10 VARCHAR(60) NOT NULL, edition_id INTEGER NOT NULL, CONSTRAINT fk_editions_isbn10_key_1 FOREIGN KEY (edition_id) REFERENCES editions (id))",
		"CREATE TABLE IF NOT EXISTS editions_isbn13 ( id INTEGER PRIMARY KEY AUTOINCREMENT, isbn13 VARCHAR(60) NOT NULL, edition_id INTEGER NOT NULL, CONSTRAINT fk_editions_isbn10_key_1 FOREIGN KEY (edition_id) REFERENCES editions (id))",
		"CREATE TABLE IF NOT EXISTS editions_authors ( id INTEGER PRIMARY KEY AUTOINCREMENT, author_id INTEGER NOT NULL, edition_id INTEGER NOT NULL, CONSTRAINT fk_authors_editions_key_1 FOREIGN KEY (author_id) REFERENCES authors (id), CONSTRAINT fk_authors_editions_key_2 FOREIGN KEY (edition_id) REFERENCES editions (id))",
	}

	tx, err := db.conn.Begin()
	if err != nil {
		panic("Cannot start transaction: " + err.Error())
	}

	for _, table := range tables {
		_, err = tx.Exec(table)
		if err != nil {
			panic("Cannot create table: " + err.Error())
		}
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
