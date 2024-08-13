package dbpg

import (
	"database/sql"
	"dbtest/common"
	"fmt"
	_ "github.com/lib/pq"
)

type Postgres struct {
	size     int
	conn     *sql.DB
	host     string
	database string
	init     bool
}

func (db *Postgres) create() {
	var err error

	tables := []string{
		"CREATE TABLE IF NOT EXISTS json_data(data JSONB)",
		"CREATE TABLE IF NOT EXISTS authors (id SERIAL PRIMARY KEY, key VARCHAR(60) NOT NULL, name VARCHAR(60) NOT NULL, revision INTEGER, birth INTEGER,death INTEGER, CONSTRAINT uk_authors_key UNIQUE (key, name))",
		"CREATE TABLE IF NOT EXISTS editions (id SERIAL PRIMARY KEY, key VARCHAR(60) NOT NULL, title VARCHAR(60) NOT NULL, subtitle VARCHAR(120), format TEXT, publish_date TEXT, edition TEXT, description TEXT, pages INTEGER, CONSTRAINT uk_editions_key UNIQUE (key, title))",
		"CREATE TABLE IF NOT EXISTS genres (id SERIAL PRIMARY KEY, name VARCHAR(60) NOT NULL, CONSTRAINT uk_genres_key UNIQUE (name))",
		"CREATE TABLE IF NOT EXISTS publishers (id SERIAL PRIMARY KEY, name VARCHAR(60) NOT NULL, CONSTRAINT uk_publishers_key UNIQUE (name))",
		"CREATE TABLE IF NOT EXISTS editions_genres (id SERIAL PRIMARY KEY, edition_id INTEGER, genre_id INTEGER, CONSTRAINT fk_editions_genres_key_1 FOREIGN KEY (edition_id) REFERENCES editions (id), CONSTRAINT fk_editions_genres_key_2 FOREIGN KEY (genre_id) REFERENCES genres (id))",
		"CREATE TABLE IF NOT EXISTS editions_publishers (id SERIAL PRIMARY KEY, edition_id INTEGER, publisher_id INTEGER, CONSTRAINT fk_editions_publishers_key_1 FOREIGN KEY (edition_id) REFERENCES editions (id), CONSTRAINT fk_editions_publishers_key_2 FOREIGN KEY (publisher_id) REFERENCES publishers (id))",
		"CREATE TABLE IF NOT EXISTS editions_isbn10 (id SERIAL PRIMARY KEY, isbn10 VARCHAR(60) NOT NULL, edition_id INTEGER NOT NULL, CONSTRAINT fk_editions_isbn10_key_1 FOREIGN KEY (edition_id) REFERENCES editions (id))",
		"CREATE TABLE IF NOT EXISTS editions_isbn13 (id SERIAL PRIMARY KEY, isbn13 VARCHAR(60) NOT NULL, edition_id INTEGER NOT NULL, CONSTRAINT fk_editions_isbn10_key_1 FOREIGN KEY (edition_id) REFERENCES editions (id))",
		"CREATE TABLE IF NOT EXISTS editions_authors (id SERIAL PRIMARY KEY, author_id INTEGER NOT NULL, edition_id INTEGER NOT NULL, CONSTRAINT fk_authors_editions_key_1 FOREIGN KEY (author_id) REFERENCES authors (id), CONSTRAINT fk_authors_editions_key_2 FOREIGN KEY (edition_id) REFERENCES editions (id))",
	}

	for _, table := range tables {
		_, err = db.conn.Exec(table)
		if err != nil {
			panic("Cannot create table: " + err.Error())
		}
	}
}

func (db *Postgres) clean() {
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

func (db *Postgres) New(cli *common.CLI) {
	var err error

	db.size = cli.Rows

	db.host = cli.Postgres.Host
	db.database = cli.Postgres.Database
	db.init = cli.Init

	port := cli.Postgres.Port
	user := cli.Postgres.User
	password := cli.Postgres.Password

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", user, password, db.database, db.host, port)
	db.conn, err = sql.Open("postgres", dsn)
	if err != nil {
		panic("Cannot open database connection: " + err.Error())
	}

	if db.init {
		db.clean()
	}
	db.create()
}

func (db *Postgres) Close() {
	db.conn.Close()
}

func (db *Postgres) Name() string {
	return "PostgreSQL"
}

func (db *Postgres) Url() string {
	return db.host
}
