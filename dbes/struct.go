package dbes

import (
	"bytes"
	"dbtest/common"
	"fmt"
	"io"
	"net/http"
)

type Elasticsearch struct {
	host     string
	port     int
	index    string
	shards   int
	replicas int
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
	mapping := fmt.Sprintf(`{
		"settings" : {
			"index" : {
				"number_of_shards" : %d,
				"number_of_replicas" : %d
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
	}`, e.shards, e.replicas)

	e.call("PUT", fmt.Sprintf("http://%s:%d/%s", e.host, e.port, e.index), bytes.NewBuffer([]byte(mapping)))
}

func (e *Elasticsearch) New(cli *common.CLI) {
	e.host = cli.Elasticsearch.Host
	e.port = cli.Elasticsearch.Port
	e.index = cli.Elasticsearch.Index
	e.shards = cli.Elasticsearch.Shards
	e.replicas = cli.Elasticsearch.Replicas

	if cli.Init {
		e.clean()
	}
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
