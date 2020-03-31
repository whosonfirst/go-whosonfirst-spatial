package extras

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strings"
)

func AppendExtrasWithBytes(ctx context.Context, source []byte, target []byte, extras []string) ([]byte, error) {

	var err error

	for _, e := range extras {

		paths := make([]string, 0)

		if strings.HasSuffix(e, "*") || strings.HasSuffix(e, ":") {

			e = strings.Replace(e, "*", "", -1)

			props := gjson.GetBytes(source, "properties")

			for k, _ := range props.Map() {

				if strings.HasPrefix(k, e) {
					paths = append(paths, k)
				}
			}

		} else {
			paths = append(paths, e)
		}

		for _, p := range paths {

			get_path := fmt.Sprintf("properties.%s", p)
			set_path := fmt.Sprintf("properties.%s", p) // FIX ME

			v := gjson.GetBytes(source, get_path)

			/*
				log.Println("GET", id, get_path)
				log.Println("SET", id, set_path)
				log.Println("VALUE", v.Value())
			*/

			if !v.Exists() {
				continue
			}

			target, err = sjson.SetBytes(target, set_path, v.Value())

			if err != nil {
				return nil, err
			}
		}
	}

	return target, nil
}
