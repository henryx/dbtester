package dbsqlite

import "log"

func (db *SQLite) Load(size int, filename string) {
	if !db.init {
		log.Println("Skipped load JSON data to database")
		return
	}

	if filename == "" {
		panic("No datafile specified")
	}

	db.loadJSON(size, filename)
	db.loadSchema()
}
