package dbsqlite

import "log"

func (db *SQLite) Load(filename string) {

	if !db.init {
		log.Println("Skipped load JSON data to database")
		return
	}

	if filename == "" {
		panic("No datafile specified")
	}

	db.loadJSON(db.size, filename)
	db.loadSchema()
}
