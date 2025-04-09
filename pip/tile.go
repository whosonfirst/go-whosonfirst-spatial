package pip

import (
	"context"
	"fmt"
	"strconv"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/clip"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/maptile"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/query"
)

func PointInPolygonCandidateFeaturessFromTile(ctx context.Context, db database.SpatialDatabase, q *query.SpatialQuery, t maptile.Tile) (*geojson.FeatureCollection, error) {

	tile_bounds := t.Bound()
	tile_geom := tile_bounds.ToPolygon()

	q.Geometry = geojson.NewGeometry(tile_geom)

	intersects_fn, err := query.NewSpatialFunction(ctx, "intersects://")

	if err != nil {
		return nil, fmt.Errorf("Failed to construct spatial fuction (intersects://), %w", err)
	}

	intersects_rsp, err := query.ExecuteQuery(ctx, db, intersects_fn, q)

	if err != nil {
		return nil, fmt.Errorf("Failed to execute query, %w", err)
	}

	// To do: For each result:
	// Fetch geojson Feature
	// Trim/clip geometries to maptile
	// Return GeoJSON

	fc := geojson.NewFeatureCollection()

	for _, r := range intersects_rsp.Results() {

		id, err := strconv.ParseInt(r.Id(), 10, 64)

		if err != nil {
			return nil, err
		}

		body, err := wof_reader.LoadBytes(ctx, db, id)

		if err != nil {
			return nil, err
		}

		f, err := geojson.UnmarshalFeature(body)

		if err != nil {
			return nil, err
		}

		// Clipping happens below
		fc.Append(f)
	}

	col := make([]orb.Geometry, len(fc.Features))

	for idx, f := range fc.Features {
		col[idx] = f.Geometry
	}

	col = clip.Collection(tile_bounds, col)

	for idx, clipped_geom := range col {
		fc.Features[idx].Geometry = clipped_geom
	}

	return fc, nil
}
