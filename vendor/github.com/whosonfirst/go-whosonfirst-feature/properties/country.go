package properties

import (
	"github.com/tidwall/gjson"
)

func Country(body []byte) string {

	rsp := gjson.GetBytes(body, "properties.wof:country")

	if !rsp.Exists() {
		return "XY"
	}

	return rsp.String()
}
