package dbcouch

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func (c *CouchDB) save(buf string) {
	log.Println("Save CouchDB buffer data:\n", buf)
	resp := c.call("POST", fmt.Sprintf("http://%s:%s@%s:%d/%s/_bulk_docs", c.user, c.password, c.host, c.port, c.database), bytes.NewBuffer([]byte(buf)))

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

func (c *CouchDB) Load(size int, filename string) {
	var err error
	var buf strings.Builder
	var data strings.Builder

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
				if buf.Len() > 0 {
					data.Reset()

					buffer := buf.String()
					buffer = strings.TrimSuffix(buffer, ",")

					data.WriteString(`{"docs": [`)
					data.WriteString(buffer)
					data.WriteString(`]}`)
					c.save(data.String())
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

		line = strings.TrimRightFunc(line, func(c rune) bool {
			//In windows newline is \r\n
			return c == '\r' || c == '\n'
		})
		buf.WriteString(line)
		buf.WriteString(",")

		if counter == size {
			data.Reset()

			buffer := buf.String()
			buffer = strings.TrimSuffix(buffer, ",")

			data.WriteString(`{"docs": [`)
			data.WriteString(buffer)
			data.WriteString(`]}`)

			c.save(data.String())
			commit++
			log.Printf("Committed %d...\n", commit)

			buf.Reset()
			counter = 0
		}
	}
}
