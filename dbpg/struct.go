package dbpg

import (
	"database/sql"
	"dbtest/cli"
	"fmt"
	_ "github.com/lib/pq"
)

type Postgres struct {
	conn     *sql.DB
	host     string
	database string
}

func (db *Postgres) create() {
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

func (db *Postgres) clean() {
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

func (db *Postgres) New(cli *cli.CLI) {
	var err error

	db.host = cli.Postgres.Host
	port := cli.Postgres.Port
	user := cli.Postgres.User
	password := cli.Postgres.Password
	db.database = cli.Postgres.Database

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
