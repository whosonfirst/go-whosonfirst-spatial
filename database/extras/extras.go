package extras

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func AppendExtrasWithBytes(ctx context.Context, source []byte, target []byte, extras []string) ([]byte, error) {

	var err error

	for _, e := range extras {

		// TO DO: CHECK FOR WILDCARDS

		path := fmt.Sprintf("properties.%s", e)
		e_rsp := gjson.GetBytes(source, path)

		if !e_rsp.Exists() {
			continue
		}

		target, err = sjson.SetBytes(target, path, e_rsp.Value())

		if err != nil {
			return nil, err
		}
	}

	return target, nil
}
