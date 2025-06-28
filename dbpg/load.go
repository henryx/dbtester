package dbpg

import (
	"log"
)

func (db *Postgres) Load(filename string) {
	if !db.init {
		log.Println("Skipped load JSON data to database")
		return
	}

	if filename == "" {
		panic("No datafile specified")
	}

	db.loadJSON(db.rows, filename)
	if db.transform {
		db.loadSchema()
	}
}
