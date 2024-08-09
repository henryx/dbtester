package dbsqlite

import (
	"database/sql"
	"dbtest/common"
	"errors"
	"github.com/tidwall/gjson"
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
	var insert = "INSERT INTO editions(key, title, subtitle, format, publish_date, edition, description, pages) VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id"

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

	description := ""
	switch desc.Type {
	case gjson.String:
		description = desc.String()
	case gjson.JSON:
		description = desc.Get("description.value").String()
	}

	pages := gjson.Get(j, "number_of_pages")

	stmt, err := tx.Prepare(insert)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}
	defer stmt.Close()

	err = stmt.QueryRow(key.String(), title.String(), nullString(subtitle), nullString(format),
		nullString(publishDate), nullString(edition), description, nullString(pages)).Scan(&editionId)

	publishersIds := db.getPublishers(tx, publishers.Array())
	db.addEditionPublishers(tx, editionId, publishersIds)

	genresIds := db.getGenres(tx, genres.Array())
	db.addEditionGenres(tx, editionId, genresIds)

	authorsIds := db.getAuthors(tx, authors.Array())
	db.addEditionAuthors(tx, editionId, authorsIds)
	// TODO: save ISBNs
}

func (db *SQLite) addEditionAuthors(tx *sql.Tx, editionId int, authors []int) {
	// TODO: add authors
}

func (db *SQLite) addEditionPublishers(tx *sql.Tx, editionId int, publishers []int) {
	var err error

	insert := "INSERT INTO editions_publishers(edition_id, publisher_id) VALUES(?, ?)"

	stmt, err := tx.Prepare(insert)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}
	defer stmt.Close()

	for _, publisher := range publishers {
		_, err := stmt.Exec(editionId, publisher)
		if err != nil {
			panic("Cannot insert publishers/editions relation: " + err.Error())
		}
	}
}

func (db *SQLite) addEditionGenres(tx *sql.Tx, editionId int, genres []int) {
	var err error

	insert := "INSERT INTO editions_genres(edition_id, genre_id) VALUES(?, ?)"

	stmt, err := tx.Prepare(insert)
	if err != nil {
		panic("Cannot create statement: " + err.Error())
	}
	defer stmt.Close()

	for _, genre := range genres {
		_, err := stmt.Exec(editionId, genre)
		if err != nil {
			panic("Cannot insert genres/editions relation: " + err.Error())
		}
	}
}

func (db *SQLite) addAuthor(tx *sql.Tx, j string) {
	var err error
	var insert = "INSERT INTO authors(key, name, revision, birth, death) VALUES(?, ?, ?, ?, ?)"

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

	_, err = tx.Exec(insert, key.String(), name.String(), revision.Int(), nullInt(birth), nullInt(death))
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
			insert := "INSERT INTO publishers(name) VALUES (?) RETURNING id"
			stmt, err := tx.Prepare(insert)
			if err != nil {
				return nil
			}
			defer stmt.Close()

			err = stmt.QueryRow(publisher.String()).Scan(&publisherId)
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
			insert := "INSERT INTO genres(name) VALUES (?) RETURNING id"
			stmt, err := tx.Prepare(insert)
			if err != nil {
				return nil
			}

			err = stmt.QueryRow(genre.String()).Scan(&genreId)
			if err != nil {
				return nil
			}

			_ = stmt.Close()
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
		db.addAuthor(tx, j)
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

func (db *SQLite) loadSchema() {
	var tx *sql.Tx
	var err error

	tx, err = db.conn.Begin()
	if err != nil {
		panic("Cannot open transaction: " + err.Error())
	}

	db.loadAuthors(tx)

	err = tx.Commit()
	if err != nil {
		panic("Cannot commit transaction: " + err.Error())
	}

	tx, err = db.conn.Begin()
	if err != nil {
		panic("Cannot open transaction: " + err.Error())
	}

	db.loadEditions(tx)

	err = tx.Commit()
	if err != nil {
		panic("Cannot commit transaction: " + err.Error())
	}
}
