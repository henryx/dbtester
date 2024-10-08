package dbes

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func (e *Elasticsearch) save(buf string) {
	resp := e.call("POST", fmt.Sprintf("http://%s:%d/%s/_bulk?refresh=wait_for", e.host, e.port, e.index), bytes.NewBuffer([]byte(buf)))

	body, err := io.ReadAll(resp.Body)
	msg := string(body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		if err != nil {
			panic(err)
		} else {
			log.Println(resp.Status)
			panic(msg)
		}
	}
}

func (e *Elasticsearch) Load(filename string) {
	var err error
	var buf strings.Builder

	inFile, err := os.Open(filename)
	if err != nil {
		panic("Cannot open file " + filename)
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)

	counter := 0
	commit := 0

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Reached end of file")
				if buf.Len() > 0 {
					e.save(buf.String())
					log.Println("Committed last data")
				}
				break
			}

			panic("Error when load data: " + err.Error())
		}

		if strings.Trim(line, "\r\n") == "" {
			log.Println("Line empty")
			continue
		}

		counter++

		buf.WriteString(`{ "index" : { } }`)
		buf.WriteString("\n")
		buf.WriteString(line)

		if counter == e.rows {
			e.save(buf.String())
			commit++
			log.Printf("Committed %d...\n", commit)

			buf.Reset()
			counter = 0
		}
	}
}
