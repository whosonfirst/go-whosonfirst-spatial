package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/paulmach/orb/geojson"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
)

func main() {

	flag.Parse()

	fc := geojson.NewFeatureCollection()

	for _, path := range flag.Args() {

		body, err := os.ReadFile(path)

		if err != nil {
			log.Fatalf("Failed to read %s, %v", path, err)
		}

		f, err := geojson.UnmarshalFeature(body)

		if err != nil {
			log.Fatalf("Failed to unmarshal %s, %v", path, err)
		}

		pt, err := geo.FindAnchorPoint(f.Geometry)

		if err != nil {
			log.Fatalf("Failed to derive anchor point for %s, %v", path, err)
		}

		f2 := geojson.NewFeature(pt)
		
		f2.Properties = map[string]any{
			"feature": path,
		}

		fc.Append(f2)
	}

	enc := json.NewEncoder(os.Stdout)
	err := enc.Encode(fc)

	if err != nil {
		log.Fatalf("Failed to encode results, %v", err)
	}
}
