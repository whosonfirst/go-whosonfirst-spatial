package main

// go run -mod vendor cmd/wof-pip/main.go -index 'rtree:///?cache=gocache://' -cache gocache:// -mode repo:// /usr/local/data/sfomuseum-data-maps/

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	geojson_utils "github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	_ "github.com/whosonfirst/go-whosonfirst-spatial/database/rtree"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	log "log"
	"os"
	"strconv"
	"strings"
)

func main() {

	fl, err := flags.CommonFlags()

	if err != nil {
		log.Fatal(err)
	}

	flags.Parse(fl)

	err = flags.ValidateCommonFlags(fl)

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pip, err := app.NewSpatialApplicationWithFlagSet(ctx, fl)

	if err != nil {
		log.Fatal("Failed to create new PIP application, because", err)
	}

	paths := fl.Args()

	err = pip.IndexPaths(ctx, paths...)

	if err != nil {
		pip.Logger.Fatal("Failed to index paths, because %s", err)
	}

	f, err := filter.NewSPRFilter()

	if err != nil {
		pip.Logger.Fatal("Failed to create SPR filter, because %s", err)
	}

	fmt.Println("ready to query")

	spatial_db := pip.SpatialDatabase

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		input := scanner.Text()
		pip.Logger.Status("# %s", input)

		parts := strings.Split(input, " ")

		if len(parts) == 0 {
			pip.Logger.Warning("Invalid input")
			continue
		}

		var command string

		switch parts[0] {

		case "candidates":
			command = parts[0]
		case "pip":
			command = parts[0]
		default:
			pip.Logger.Warning("Invalid command")
			continue
		}

		var results interface{}

		if command == "pip" || command == "candidates" {

			str_lat := strings.Trim(parts[1], " ")
			str_lon := strings.Trim(parts[2], " ")

			lat, err := strconv.ParseFloat(str_lat, 64)

			if err != nil {
				pip.Logger.Warning("Invalid latitude, %s", err)
				continue
			}

			lon, err := strconv.ParseFloat(str_lon, 64)

			if err != nil {
				pip.Logger.Warning("Invalid longitude, %s", err)
				continue
			}

			c, err := geojson_utils.NewCoordinateFromLatLons(lat, lon)

			if err != nil {
				pip.Logger.Warning("Invalid latitude, longitude, %s", err)
				continue
			}

			if command == "pip" {

				intersects, err := spatial_db.PointInPolygon(ctx, &c, f)

				if err != nil {
					pip.Logger.Warning("Unable to get intersects, because %s", err)
					continue
				}

				results = intersects

			} else {

				candidates, err := spatial_db.PointInPolygonCandidates(ctx, &c)

				if err != nil {
					pip.Logger.Warning("Unable to get candidates, because %s", err)
					continue
				}

				results = candidates
			}

		} else {
			pip.Logger.Warning("Invalid command")
			continue
		}

		body, err := json.Marshal(results)

		if err != nil {
			pip.Logger.Warning("Failed to marshal results, because %s", err)
			continue
		}

		fmt.Println(string(body))
	}

	os.Exit(0)
}
