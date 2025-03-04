package database

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/fixtures"
	"github.com/whosonfirst/go-whosonfirst-spatial/fixtures/microhoods"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
)

type PointInPolygonCriteria struct {
	IsCurrent int64
	Latitude  float64
	Longitude float64
}

func TestSpatialDatabasePointInPolygon(t *testing.T) {

	ctx := context.Background()

	database_uri := "rtree://"

	tests := map[int64]PointInPolygonCriteria{
		1108712253: PointInPolygonCriteria{Longitude: -71.120168, Latitude: 42.376015, IsCurrent: 1},   // Old Cambridge
		420561633:  PointInPolygonCriteria{Longitude: -122.395268, Latitude: 37.794893, IsCurrent: 0},  // Superbowl City
		420780729:  PointInPolygonCriteria{Longitude: -122.421529, Latitude: 37.743168, IsCurrent: -1}, // Liminal Zone of Deliciousness
	}

	db, err := NewSpatialDatabase(ctx, database_uri)

	if err != nil {
		t.Fatalf("Failed to create new spatial database, %v", err)
	}

	path_microhoods, err := filepath.Abs("../fixtures/microhoods")

	if err != nil {
		t.Fatalf("Failed to derive path for microhoods, %v", err)
	}

	err = IndexDatabaseWithIterator(ctx, db, "directory://", path_microhoods)

	if err != nil {
		t.Fatalf("Failed to index spatial database, %v", err)
	}

	for expected, criteria := range tests {

		c, err := geo.NewCoordinate(criteria.Longitude, criteria.Latitude)

		if err != nil {
			t.Fatalf("Failed to create new coordinate, %v", err)
		}

		i, err := filter.NewSPRInputs()

		if err != nil {
			t.Fatalf("Failed to create SPR inputs, %v", err)
		}

		i.IsCurrent = []int64{criteria.IsCurrent}
		// i.Placetypes = []string{"microhood"}

		f, err := filter.NewSPRFilterFromInputs(i)

		if err != nil {
			t.Fatalf("Failed to create SPR filter from inputs, %v", err)
		}

		spr, err := db.PointInPolygon(ctx, c, f)

		if err != nil {
			t.Fatalf("Failed to perform point in polygon query, %v", err)
		}

		results := spr.Results()
		count := len(results)

		if count != 1 {
			t.Fatalf("Expected 1 result but got %d for '%d'", count, expected)
		}

		first := results[0]

		if first.Id() != strconv.FormatInt(expected, 10) {
			t.Fatalf("Expected %d but got %s", expected, first.Id())
		}
	}
}

// This is known to fail until we keep a local lookup table of all the bounding boxes associated
// with a feature is created. The way we're doing things in database.RemoveFeature using a comparator
// doesn't actually work...

func TestSpatialDatabaseRemoveFeature(t *testing.T) {

	t.Skip()

	ctx := context.Background()

	database_uri := "rtree://"

	db, err := NewSpatialDatabase(ctx, database_uri)

	if err != nil {
		t.Fatalf("Failed to create new spatial database, %v", err)
	}

	defer db.Close(ctx)

	id := 101737491
	lat := 46.852675
	lon := -71.330873

	test_data := fmt.Sprintf("fixtures/%d.geojson", id)

	fh, err := os.Open(test_data)

	if err != nil {
		t.Fatalf("Failed to open %s, %v", test_data, err)
	}

	defer fh.Close()

	body, err := io.ReadAll(fh)

	if err != nil {
		t.Fatalf("Failed to read %s, %v", test_data, err)
	}

	err = db.IndexFeature(ctx, body)

	if err != nil {
		t.Fatalf("Failed to index %s, %v", test_data, err)
	}

	c, err := geo.NewCoordinate(lon, lat)

	if err != nil {
		t.Fatalf("Failed to create new coordinate, %v", err)
	}

	spr, err := db.PointInPolygon(ctx, c)

	if err != nil {
		t.Fatalf("Failed to perform point in polygon query, %v", err)
	}

	results := spr.Results()
	count := len(results)

	if count != 1 {
		t.Fatalf("Expected 1 result but got %d", count)
	}

	err = db.RemoveFeature(ctx, "101737491")

	if err != nil {
		t.Fatalf("Failed to remove %s, %v", test_data, err)
	}

	spr, err = db.PointInPolygon(ctx, c)

	if err != nil {
		t.Fatalf("Failed to perform point in polygon query, %v", err)
	}

	results = spr.Results()
	count = len(results)

	if count != 0 {
		t.Fatalf("Expected 0 results but got %d", count)
	}
}

func TestSpatialDatabaseWithFS(t *testing.T) {

	ctx := context.Background()

	database_uri := "rtree://?dsn=:memory:"

	tests := map[int64]PointInPolygonCriteria{
		1108712253: PointInPolygonCriteria{Longitude: -71.120168, Latitude: 42.376015, IsCurrent: 1},   // Old Cambridge
		420561633:  PointInPolygonCriteria{Longitude: -122.395268, Latitude: 37.794893, IsCurrent: 0},  // Superbowl City
		420780729:  PointInPolygonCriteria{Longitude: -122.421529, Latitude: 37.743168, IsCurrent: -1}, // Liminal Zone of Deliciousness
	}

	db, err := NewSpatialDatabase(ctx, database_uri)

	if err != nil {
		t.Fatalf("Failed to create new spatial database, %v", err)
	}

	err = IndexDatabaseWithFS(ctx, db, microhoods.FS)

	if err != nil {
		t.Fatalf("Failed to index spatial database, %v", err)
	}

	for expected, criteria := range tests {

		c, err := geo.NewCoordinate(criteria.Longitude, criteria.Latitude)

		if err != nil {
			t.Fatalf("Failed to create new coordinate, %v", err)
		}

		i, err := filter.NewSPRInputs()

		if err != nil {
			t.Fatalf("Failed to create SPR inputs, %v", err)
		}

		i.IsCurrent = []int64{criteria.IsCurrent}
		// i.Placetypes = []string{"microhood"}

		f, err := filter.NewSPRFilterFromInputs(i)

		if err != nil {
			t.Fatalf("Failed to create SPR filter from inputs, %v", err)
		}

		spr, err := db.PointInPolygon(ctx, c, f)

		if err != nil {
			t.Fatalf("Failed to perform point in polygon query, %v", err)
		}

		results := spr.Results()
		count := len(results)

		if count != 1 {
			t.Fatalf("Expected 1 result but got %d for '%d'", count, expected)
		}

		first := results[0]

		if first.Id() != strconv.FormatInt(expected, 10) {
			t.Fatalf("Expected %d but got %s", expected, first.Id())
		}
	}
}

func TestSpatialDatabaseWithFeatureCollection(t *testing.T) {

	ctx := context.Background()

	database_uri := "rtree://?dsn=:memory:"

	tests := map[int64]PointInPolygonCriteria{
		1108712253: PointInPolygonCriteria{Longitude: -71.120168, Latitude: 42.376015, IsCurrent: 1},   // Old Cambridge
		420561633:  PointInPolygonCriteria{Longitude: -122.395268, Latitude: 37.794893, IsCurrent: 0},  // Superbowl City
		420780729:  PointInPolygonCriteria{Longitude: -122.421529, Latitude: 37.743168, IsCurrent: -1}, // Liminal Zone of Deliciousness
	}

	db, err := NewSpatialDatabase(ctx, database_uri)

	if err != nil {
		t.Fatalf("Failed to create new spatial database, %v", err)
	}

	err = IndexDatabaseWithFS(ctx, db, fixtures.FS)

	if err != nil {
		t.Fatalf("Failed to index spatial database, %v", err)
	}

	for expected, criteria := range tests {

		c, err := geo.NewCoordinate(criteria.Longitude, criteria.Latitude)

		if err != nil {
			t.Fatalf("Failed to create new coordinate, %v", err)
		}

		i, err := filter.NewSPRInputs()

		if err != nil {
			t.Fatalf("Failed to create SPR inputs, %v", err)
		}

		i.IsCurrent = []int64{criteria.IsCurrent}
		// i.Placetypes = []string{"microhood"}

		f, err := filter.NewSPRFilterFromInputs(i)

		if err != nil {
			t.Fatalf("Failed to create SPR filter from inputs, %v", err)
		}

		spr, err := db.PointInPolygon(ctx, c, f)

		if err != nil {
			t.Fatalf("Failed to perform point in polygon query, %v", err)
		}

		results := spr.Results()
		count := len(results)

		if count != 1 {
			t.Fatalf("Expected 1 result but got %d for '%d'", count, expected)
		}

		first := results[0]

		if first.Id() != strconv.FormatInt(expected, 10) {
			t.Fatalf("Expected %d but got %s", expected, first.Id())
		}
	}
}
