package dbes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
)

func (e *Elasticsearch) Count() int64 {
	var result map[string]interface{}

	resp := e.call("GET", fmt.Sprintf("%s/%s/_count", e.url, e.index), nil)
	defer resp.Body.Close()

	count, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(count, &result)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
	conv := math.Round(result["count"].(float64))
	return int64(conv)
}
