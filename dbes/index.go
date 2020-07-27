package dbes

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

func (e *Elasticsearch) Index() {
	mapping := `
	{"properties": {"key": {"type": "text"}}}}`
	
	resp := e.call("PUT", fmt.Sprintf("%s/%s/_mapping", e.url, e.index), bytes.NewBuffer([]byte(mapping)))
	body, err := ioutil.ReadAll(resp.Body)
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