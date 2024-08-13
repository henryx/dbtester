package dbpg

import (
	"bufio"
	"database/sql"
	"io"
	"log"
	"os"
)

func (db *Postgres) loadJSON(size int, filename string) {
	var tx *sql.Tx
	var insert string = "INSERT INTO json_data VALUES ($1)"
	var err error

	inFile, err := os.Open(filename)
	if err != nil {
		panic("Cannot open file " + filename)
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)

	counter := 0
	commit := 0
	tx, err = db.conn.Begin()
	if err != nil {
		panic("Cannot start transaction " + err.Error())
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			if line == "" {
				log.Println("Line empty")
				continue
			}

			panic("Error when load data: " + err.Error())
		}

		_, err = tx.Exec(insert, line)
		if err != nil {
			panic(err.Error())
		}

		counter++
		if counter == size {
			counter = 0
			err = tx.Commit()
			if err != nil {
				panic("Cannot commit transaction: " + err.Error())
			}
			commit++
			log.Printf("Committed %d...\n", commit)

			tx, err = db.conn.Begin()
			if err != nil {
				panic("Cannot start transaction " + err.Error())
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		panic("Cannot commit transaction: " + err.Error())
	}
}