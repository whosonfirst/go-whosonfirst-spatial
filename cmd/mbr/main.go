package main

// > go run cmd/mbr/main.go -id 102087579 -id 102086959 -id 102085387
// -123.173825,37.053858,-121.469214,37.929824

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "github.com/whosonfirst/go-reader-http"

	"github.com/paulmach/orb"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-feature/geometry"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
)

func main() {

	var reader_uri string
	var ids multi.MultiInt64

	fs := flagset.NewFlagSet("flags")
	fs.StringVar(&reader_uri, "reader-uri", "http://data.whosonfirst.org", "A registered whosonfirst/go-reader.Reader URI.")
	fs.Var(&ids, "id", "One or more Who's On First IDs.")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Derive one or more intersecting minimum-bounding-rectangles (MBR) for one or more Who's On First IDs.\n")
		fmt.Fprintf(os.Stderr, "Usage:\n\t %s [options]\n", os.Args[0])
		fs.PrintDefaults()
	}

	flagset.Parse(fs)

	ctx := context.Background()

	r, err := reader.NewReader(ctx, reader_uri)

	if err != nil {
		log.Fatalf("Failed to create reader, %v", err)
	}

	bounds := make([]orb.Bound, 0)

	for idx, id := range ids {

		body, err := wof_reader.LoadBytes(ctx, r, id)

		if err != nil {
			log.Fatalf("Failed to load record for %d, %v", id, err)
		}

		geojson_geom, err := geometry.Geometry(body)

		if err != nil {
			log.Fatalf("Failed to derive geometry for %d, %v", id, err)
		}

		orb_geom := geojson_geom.Geometry()
		this_b := orb_geom.Bound()

		switch idx {
		case 0:
			bounds = append(bounds, this_b)
		default:

			has_rel := false

			for idx_b, other_b := range bounds {

				if this_b.Intersects(other_b) {
					bounds[idx_b] = other_b.Union(this_b)
					has_rel = true
					break
				}
			}

			if !has_rel {
				bounds = append(bounds, this_b)
			}
		}
	}

	for _, b := range bounds {
		fmt.Printf("%f,%f,%f,%f\n", b.Min.X(), b.Min.Y(), b.Max.X(), b.Max.Y())
	}
}
