package dbes

import (
	"bytes"
	"fmt"
	"io"
)

func (e *Elasticsearch) IndexJSON() {
	mapping := `
	{"properties": {"key": {"type": "text"}}}}`

	resp := e.call("PUT", fmt.Sprintf("http://%s:%d/%s/_mapping", e.host, e.port, e.index), bytes.NewBuffer([]byte(mapping)))
	body, err := io.ReadAll(resp.Body)
	msg := string(body)
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		if err != nil {
			panic(err)
		} else {
			panic(msg)
		}
	}
}
