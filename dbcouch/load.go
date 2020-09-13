package dbcouch

import (
	"bufio"
	"io"
	"log"
	"os"
)

func (e *CouchDB) save(buf string) {
	// TODO: implement save
}

func (c *CouchDB) Load(size int, filename string) {
	var err error

	inFile, err := os.Open(filename)
	if err != nil {
		panic("Cannot open file " + filename)
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)

	counter := 0
	commit := 0
	buf := ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				if buf != "" {
					c.save(buf)
				}
				break
			}

			if line == "" {
				log.Println("Line empty")
				continue
			}
			panic("Error when load data: " + err.Error())
		}
		counter++

		if counter == size {
			c.save(buf)
			commit++
			log.Printf("Committed %d...\n", commit)

			buf = ""
			counter = 0
		}
	}
}
