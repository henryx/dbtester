package dbsqlite

import (
	"database/sql"
	"dbtest/common"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type SQLite struct {
	size                  int
	conn                  *sql.DB
	database              string
	init                  bool
	stmtAuthors           *sql.Stmt
	stmtEditions          *sql.Stmt
	stmtPublishers        *sql.Stmt
	stmtGenres            *sql.Stmt
	stmtIsbn10            *sql.Stmt
	stmtIsbn13            *sql.Stmt
	stmtEditionAuthors    *sql.Stmt
	stmtEditionPublishers *sql.Stmt
	stmtEditionGenres     *sql.Stmt
}

func (db *SQLite) create() {
	tables := []string{
		"CREATE TABLE IF NOT EXISTS json_data(data JSONB)",
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

func (db *SQLite) clean() {
	tables := []string{
		"DROP TABLE IF EXISTS json_data",
		"DROP TABLE IF EXISTS editions_authors",
		"DROP TABLE IF EXISTS editions_isbn13",
		"DROP TABLE IF EXISTS editions_isbn10",
		"DROP TABLE IF EXISTS editions_publishers",
		"DROP TABLE IF EXISTS editions_genres",
		"DROP TABLE IF EXISTS publishers",
		"DROP TABLE IF EXISTS genres",
		"DROP TABLE IF EXISTS editions",
		"DROP TABLE IF EXISTS authors",
	}

	tx, err := db.conn.Begin()
	if err != nil {
		panic("Cannot start transaction: " + err.Error())
	}

	for _, table := range tables {
		_, err = tx.Exec(table)
		if err != nil {
			panic("Cannot drop table: " + err.Error())
		}
	}
	_ = tx.Commit()
}

func (db *SQLite) New(cli *common.CLI) {
	var err error

	db.database = cli.SQLite.Database
	db.init = cli.Init
	db.size = cli.Rows

	dsn := fmt.Sprintf("file:%s?_journal=WAL&_fk=true", db.database)
	db.conn, err = sql.Open("sqlite3", dsn)

	if err != nil {
		panic("Cannot open database connection: " + err.Error())
	}

	if db.init {
		db.clean()
	}
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
