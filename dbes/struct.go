package dbes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Elasticsearch struct {
	url   string
	index string
}

func (e *Elasticsearch) call(method, uri string, buffer io.Reader) *http.Response {
	client := &http.Client{}
	
	req, err := http.NewRequest(method, uri, buffer)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-ndjson")
	
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func (e *Elasticsearch) clean() {
	e.call("DELETE", fmt.Sprintf("%s/%s", e.url, e.index), bytes.NewBuffer([]byte(nil)))
}

func (e *Elasticsearch) create() {
	mapping := `{
		"settings" : {
			"index" : {
				"number_of_shards" : 1,
				"number_of_replicas" : 0
			}
		},
		"mappings": {
			"properties": {
				"authors":{"enabled": false, "type": "object"},
				"bio": {"enabled": false, "type": "object"},
				"description": {"enabled": false, "type": "object"},
				"first_sentence": {"enabled": false, "type": "object"},
				"notes": {"enabled": false, "type": "object"},
				"reason": {"enabled": false, "type": "object"},
				"table_of_contents": {"enabled": false, "type": "object"}
			}
		}
	}`
	e.call("PUT", fmt.Sprintf("%s/%s", e.url, e.index), bytes.NewBuffer([]byte(mapping)))
}

func (e *Elasticsearch) New(host string) {
	e.url = fmt.Sprintf("http://%s:9200", host)
	e.index = "libraries"

	e.clean()
	e.create()
}

func (e *Elasticsearch) Close() {
	return
}

func (e *Elasticsearch) Name() string {
	return "Elasticsearch"
}
