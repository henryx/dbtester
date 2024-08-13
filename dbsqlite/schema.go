package dbsqlite

import (
	"database/sql"
	"dbtest/common"
	"errors"
	"github.com/tidwall/gjson"
	"log"
)

func nullInt(v gjson.Result) sql.NullInt64 {
	if v.Type == gjson.Null {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: v.Int(),
		Valid: true,
	}
}

func nullString(v gjson.Result) sql.NullString {
	if v.Type == gjson.Null {
		return sql.NullString{}
	}

	return sql.NullString{
		String: v.String(),
		Valid:  true,
	}
}

func (db *SQLite) addEdition(tx *sql.Tx, j string) {
	var editionId int
	var err error

	key := gjson.Get(j, "key")
	if !key.Exists() {
		return
	}

	title := gjson.Get(j, "title")
	if !title.Exists() {
		return
	}

	subtitle := gjson.Get(j, "subtitle")
	format := gjson.Get(j, "format")
	publishDate := gjson.Get(j, "publish_date")
	edition := gjson.Get(j, "edition_name")
	publishers := gjson.Get(j, "publishers")
	genres := gjson.Get(j, "genres")
	authors := gjson.Get(j, "authors")
	desc := gjson.Get(j, "description")
	isbn10s := gjson.Get(j, "isbn_10")
	isbn13s := gjson.Get(j, "isbn_13")

	description := ""
	switch desc.Type {
	case gjson.String:
		description = desc.String()
	case gjson.JSON:
		description = desc.Get("description.value").String()
	}

	pages := gjson.Get(j, "number_of_pages")

	err = db.stmtEditions.QueryRow(key.String(), title.String(), nullString(subtitle), nullString(format),
		nullString(publishDate), nullString(edition), description, nullString(pages)).Scan(&editionId)
	if err != nil {
		log.Println(err)
	}

	publishersIds := db.getPublishers(tx, publishers.Array())
	db.addEditionPublishers(editionId, publishersIds)

	genresIds := db.getGenres(tx, genres.Array())
	db.addEditionGenres(editionId, genresIds)

	authorsIds := db.getAuthors(tx, authors.Array())
	db.addEditionAuthors(editionId, authorsIds)

	db.addISBNs(editionId, isbn10s.Array(), isbn13s.Array())
}

func (db SQLite) addISBNs(editionId int, isbn10s []gjson.Result, isbn13s []gjson.Result) {
	var err error

	for _, isbn10 := range isbn10s {
		_, err = db.stmtIsbn10.Exec(editionId, isbn10.String())
		if err != nil {
			panic("Cannot insert ISBN10/editions relation: " + err.Error())
		}
	}

	for _, isbn13 := range isbn13s {
		_, err = db.stmtIsbn13.Exec(editionId, isbn13.String())
		if err != nil {
			panic("Cannot insert ISBN13/editions relation: " + err.Error())
		}
	}
}

func (db *SQLite) addEditionAuthors(editionId int, authors []int) {
	var err error

	for _, author := range authors {
		_, err = db.stmtEditionAuthors.Exec(editionId, author)
		if err != nil {
			panic("Cannot insert authors/editions relation: " + err.Error())
		}
	}
}

func (db *SQLite) addEditionPublishers(editionId int, publishers []int) {
	var err error

	for _, publisher := range publishers {
		_, err = db.stmtEditionPublishers.Exec(editionId, publisher)
		if err != nil {
			panic("Cannot insert publishers/editions relation: " + err.Error())
		}
	}
}

func (db *SQLite) addEditionGenres(editionId int, genres []int) {
	var err error

	for _, genre := range genres {
		_, err = db.stmtEditionGenres.Exec(editionId, genre)
		if err != nil {
			panic("Cannot insert genres/editions relation: " + err.Error())
		}
	}
}

func (db *SQLite) addAuthor(j string) {
	var err error

	key := gjson.Get(j, "key")
	if !key.Exists() {
		return
	}

	name := gjson.Get(j, "name")
	if !name.Exists() {
		return
	}

	revision := gjson.Get(j, "revision")
	birth := gjson.Get(j, "birth_date")
	death := gjson.Get(j, "death_date")

	_, err = db.stmtAuthors.Exec(key.String(), name.String(), revision.Int(), nullInt(birth), nullInt(death))
	if err != nil {
		panic("Cannot execute query: " + err.Error())
	}
}

func (db *SQLite) getAuthors(tx *sql.Tx, authors []gjson.Result) []int {
	var err error
	var authorId int

	res := make([]int, 0)
	query := "SELECT id FROM authors WHERE key = ?"

	for _, author := range authors {
		var key string
		switch author.Type {
		case gjson.Null:
			continue
		case gjson.String:
			key = author.String()
		case gjson.JSON:
			key = gjson.Get(author.Raw, "key").String()
		}

		err = tx.QueryRow(query, key).Scan(&authorId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			panic("Cannot get author: " + err.Error())
		}

		res = append(res, authorId)
	}
	return res
}

func (db *SQLite) getPublishers(tx *sql.Tx, publishers []gjson.Result) []int {
	var publisherId int

	res := make([]int, 0)
	query := "SELECT id FROM publishers WHERE name = ?"

	for _, publisher := range publishers {
		err := tx.QueryRow(query, publisher.String()).Scan(&publisherId)
		if !errors.Is(err, sql.ErrNoRows) {
			if err != nil {
				panic("Cannot get publishers: " + err.Error())
			}
		}

		if publisherId == 0 {
			err = db.stmtPublishers.QueryRow(publisher.String()).Scan(&publisherId)
			if err != nil {
				return nil
			}
		}
		res = append(res, publisherId)
	}
	return res
}

func (db *SQLite) getGenres(tx *sql.Tx, genres []gjson.Result) []int {
	var genreId int

	res := make([]int, 0)
	query := "SELECT id FROM genres WHERE name = ?"

	for _, genre := range genres {
		err := tx.QueryRow(query, genre.String()).Scan(&genreId)
		if !errors.Is(err, sql.ErrNoRows) {
			if err != nil {
				panic("Cannot get genres: " + err.Error())
			}
		}

		if genreId == 0 {
			err = db.stmtGenres.QueryRow(genre.String()).Scan(&genreId)
			if err != nil {
				return nil
			}
		}
		res = append(res, genreId)
	}
	return res
}

func (db *SQLite) loadAuthors(tx *sql.Tx) {
	var err error

	query := `SELECT data FROM json_data WHERE data->'type'->>'key' = ?`
	j := ""

	row, err := tx.Query(query, common.AUTHORS)
	if err != nil {
		panic("Cannot execute query: " + err.Error())
	}

	for row.Next() {
		err = row.Scan(&j)
		if err != nil {
			panic("Cannot read value: " + err.Error())
		}
		db.addAuthor(j)
	}
}

func (db *SQLite) loadEditions(tx *sql.Tx) {
	var err error

	query := `SELECT data FROM json_data WHERE data->'type'->>'key' = ?`
	j := ""

	tx, err = db.conn.Begin()
	if err != nil {
		panic("Cannot open transaction: " + err.Error())
	}

	row, err := tx.Query(query, common.EDITIONS)
	if err != nil {
		panic("Cannot execute query: " + err.Error())
	}

	for row.Next() {
		err = row.Scan(&j)
		if err != nil {
			panic("Cannot read value: " + err.Error())
		}
		db.addEdition(tx, j)
	}
}

func (db *SQLite) prepareStatements(tx *sql.Tx) {
	var err error

	insAuthors := "INSERT INTO authors(key, name, revision, birth, death) VALUES(?, ?, ?, ?, ?)"
	insEditions := "INSERT INTO editions(key, title, subtitle, format, publish_date, edition, description, pages) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id"
	insPublishers := "INSERT INTO publishers(name) VALUES (?) RETURNING id"
	insGenres := "INSERT INTO genres(name) VALUES (?) RETURNING id"
	insIsbn10 := "INSERT INTO editions_isbn10 (edition_id, isbn10) VALUES (?, ?)"
	insIsbn13 := "INSERT INTO editions_isbn13 (edition_id, isbn13) VALUES (?, ?)"
	insEditionAuthors := "INSERT INTO editions_authors(edition_id, author_id) VALUES(?, ?)"
	insEditionPublishers := "INSERT INTO editions_publishers(edition_id, publisher_id) VALUES(?, ?)"
	insEditionGenres := "INSERT INTO editions_genres(edition_id, genre_id) VALUES(?, ?)"

	db.stmtAuthors, err = tx.Prepare(insAuthors)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}

	db.stmtEditions, err = tx.Prepare(insEditions)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}

	db.stmtPublishers, err = tx.Prepare(insPublishers)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}

	db.stmtGenres, err = tx.Prepare(insGenres)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}

	db.stmtIsbn10, err = tx.Prepare(insIsbn10)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}

	db.stmtIsbn13, err = tx.Prepare(insIsbn13)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}

	db.stmtEditionAuthors, err = tx.Prepare(insEditionAuthors)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}

	db.stmtEditionPublishers, err = tx.Prepare(insEditionPublishers)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}

	db.stmtEditionGenres, err = tx.Prepare(insEditionGenres)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}
}

func (db *SQLite) closeStatements() {
	_ = db.stmtAuthors.Close()
	_ = db.stmtEditions.Close()
	_ = db.stmtPublishers.Close()
	_ = db.stmtGenres.Close()
	_ = db.stmtIsbn10.Close()
	_ = db.stmtIsbn13.Close()
	_ = db.stmtEditionAuthors.Close()
	_ = db.stmtEditionPublishers.Close()
}

func (db *SQLite) loadSchema() {
	var tx *sql.Tx
	var err error

	tx, err = db.conn.Begin()
	if err != nil {
		panic("Cannot open transaction: " + err.Error())
	}

	db.prepareStatements(tx)

	log.Println("Loading authors table...")
	db.loadAuthors(tx)

	log.Println("Loading editions table...")
	db.loadEditions(tx)

	db.closeStatements()
	err = tx.Commit()
	if err != nil {
		panic("Cannot commit transaction: " + err.Error())
	}
}
