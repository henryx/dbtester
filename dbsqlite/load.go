package dbsqlite

import (
	"bufio"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

type Items map[string]interface{}

func (a *Items) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Items) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func (db *SQLite) Load(size int, filename string) {
	/* j := Items{} */
	var tx *sql.Tx
	var insert string = "INSERT INTO json_data VALUES (?)"
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

		/* 		err = j.Scan([]byte(line))
		   		if err != nil {
		   			panic("Error when unmarshal data: " + err.Error())
		   		} */

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
