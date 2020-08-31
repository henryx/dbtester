package dbes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type result struct {
	Hits hits `json:"hits"`
}
type hits struct {
	Total total `json:"total"`
}
type total struct {
	Value float64 `json:"value"`
}

func (e *Elasticsearch) Find() int64 {
	var res result

	query := `{"query": {"term": {"key.keyword": "/books/OL17806216M"}}}`

	resp := e.call("GET", fmt.Sprintf("%s:%d/%s/_search", e.host, e.port, e.index), bytes.NewBuffer([]byte(query)))
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(body, &res)
	hits := res.Hits.Total.Value

	return int64(hits)
}
