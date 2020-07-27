package dbpg

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Postgres struct {
	conn  *sql.DB
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

func (db *Postgres) New(host string) {
	var err error

	dsn := strings.Join([]string{
		"user=" + "postgres",
		"password=" + "Latina,1",
		"dbname=" + "postgres",
		fmt.Sprintf("host=%s", host),
		"sslmode=disable",
	}, " ")
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
