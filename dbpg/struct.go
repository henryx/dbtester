package dbpg

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"gopkg.in/ini.v1"
)

type Postgres struct {
	conn     *sql.DB
	host     string
	database string
}

func (db *Postgres) create() {
	query := "CREATE TABLE IF NOT EXISTS test(data JSONB)"

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

func (db *Postgres) clean() {
	query := "DROP TABLE IF EXISTS test"

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

func (db *Postgres) New(cfg *ini.Section) {
	var err error

	db.host = cfg.Key("host").MustString("localhost")
	port := cfg.Key("port").MustInt(5432)
	user := cfg.Key("user").MustString("postgres")
	password := cfg.Key("password").MustString("postgres")
	db.database = cfg.Key("database").MustString("postgres")

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable", user, password, db.database, db.host, port)
	db.conn, err = sql.Open("postgres", dsn)
	if err != nil {
		panic("Cannot open database connection: " + err.Error())
	}

	db.clean()
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
