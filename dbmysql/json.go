package dbmysql

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func (db *MySQL) loadJSON(size int, filename string) {
	/* j := Items{} */
	var tx *sql.Tx
	var insert string = fmt.Sprintf("INSERT INTO %s.json_data(data) VALUES (?)", db.database)
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

			panic("Error when load data: " + err.Error())
		}

		if strings.Trim(line, "\r\n") == "" {
			log.Println("Line empty")
			continue
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
		if counter == db.rows {
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
