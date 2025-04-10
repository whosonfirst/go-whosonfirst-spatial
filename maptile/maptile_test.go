package maptile

import (
	"context"
	_ "fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/paulmach/orb/maptile"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/query"
)

func TestPointInPolygonCandidateFeaturessFromTile(t *testing.T) {

	ctx := context.Background()

	path_microhoods, err := filepath.Abs("../fixtures/microhoods")

	if err != nil {
		t.Fatalf("Failed to derive path for microhoods, %v", err)
	}

	database_uri := "rtree://"

	db, err := database.NewSpatialDatabase(ctx, database_uri)

	if err != nil {
		t.Fatalf("Failed to create new spatial database, %v", err)
	}

	err = database.IndexDatabaseWithIterator(ctx, db, "directory://?_exclude=.*\\.go&_exclude=.*~", path_microhoods)

	if err != nil {
		t.Fatalf("Failed to index spatial database, %v", err)
	}

	// As in: https://tile.openstreetmap.org/16/10482/25328.png

	tests := map[string]int{
		"16/10482/25328": 11,
		"13/1308/3165":   1,
		"17/20968/50656": 2,
	}

	for path, expected_count := range tests {

		parts := strings.Split(path, "/")

		z, err := strconv.Atoi(parts[0])

		if err != nil {
			t.Fatalf("Failed to parse '%s' (%s), %v", parts[0], path, err)
		}

		x, err := strconv.Atoi(parts[1])

		if err != nil {
			t.Fatalf("Failed to parse '%s' (%s), %v", parts[1], path, err)
		}

		y, err := strconv.Atoi(parts[2])

		if err != nil {
			t.Fatalf("Failed to parse '%s' (%s), %v", parts[2], path, err)
		}

		zm := maptile.Zoom(uint32(z))
		map_t := maptile.New(uint32(x), uint32(y), zm)

		spatial_q := &query.SpatialQuery{}

		fc, err := PointInPolygonCandidateFeaturesFromTile(ctx, db, spatial_q, &map_t)

		if err != nil {
			t.Fatalf("Failed to derive feature collection from tile query, %v", err)
		}

		count := len(fc.Features)

		if count != expected_count {
			t.Fatalf("Expected %d results for %s but got %d", expected_count, path, count)
		}
	}
}
