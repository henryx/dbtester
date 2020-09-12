package dbes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/ini.v1"
)

type Elasticsearch struct {
	host  string
	port  int
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
	e.call("DELETE", fmt.Sprintf("http://%s:%d/%s", e.host, e.port, e.index), bytes.NewBuffer([]byte(nil)))
}

func (e *Elasticsearch) create() {
	mapping := `{
		"settings" : {
			"index" : {
				"number_of_shards" : 5,
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
	e.call("PUT", fmt.Sprintf("http://%s:%d/%s", e.host, e.port, e.index), bytes.NewBuffer([]byte(mapping)))
}

func (e *Elasticsearch) New(cfg *ini.Section) {
	e.host = cfg.Key("host").MustString("localhost")
	e.port = cfg.Key("port").MustInt(9200)
	e.index = cfg.Key("index").MustString("test")

	e.clean()
	e.create()
}

func (e *Elasticsearch) Close() {
	return
}

func (e *Elasticsearch) Name() string {
	return "Elasticsearch"
}

func (e *Elasticsearch) Url() string {
	return e.host
}
