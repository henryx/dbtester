package dbes

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func (e *Elasticsearch) save(buf string) {
	resp := e.call("POST", fmt.Sprintf("%s/%s/_bulk?refresh=wait_for", e.url, e.index), bytes.NewBuffer([]byte(buf)))

	body, err := ioutil.ReadAll(resp.Body)
	msg := string(body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		if err != nil {
			panic(err)
		} else {
			log.Println(buf)
			panic(msg)
		}
	}
}

func (e *Elasticsearch) Load(size int, filename string) {
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
					e.save(buf)
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

		buf = buf + `{ "index" : { } }` + "\n" + line
		if counter == size {
			e.save(buf)
			commit++
			log.Printf("Committed %d...\n", commit)

			buf = ""
			counter = 0
		}
	}
}