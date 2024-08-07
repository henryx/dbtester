package dbes

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
)

func (e *Elasticsearch) CountJSON() int64 {
	var result map[string]interface{}

	resp := e.call("GET", fmt.Sprintf("http://%s:%d/%s/_count", e.host, e.port, e.index), nil)
	defer resp.Body.Close()

	count, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(count, &result)

	conv := math.Round(result["count"].(float64))
	return int64(conv)
}
