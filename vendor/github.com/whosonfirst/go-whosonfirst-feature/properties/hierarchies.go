package properties

import (
	"github.com/tidwall/gjson"
)

func Hierarchies(body []byte) []map[string]int64 {

	rsp := gjson.GetBytes(body, "properties.wof:hierarchy")

	if !rsp.Exists() {
		return nil
	}

	hierarchies := make([]map[string]int64, 0)

	for _, h := range rsp.Array() {

		dict := make(map[string]int64)

		for k, v := range h.Map() {
			dict[k] = v.Int()
		}

		hierarchies = append(hierarchies, dict)
	}

	return hierarchies
}
