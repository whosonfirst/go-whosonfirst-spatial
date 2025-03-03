package request

import (
	"context"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type IntersectsQuery struct {
	Query
}

func NewIntersectsQuery(ctx context.Context, uri string) (Query, error) {
	q := &IntersectsQuery{}
	return q, nil
}

func (q *IntersectsQuery) Execute(ctx context.Context, db database.SpatialDatabase, geom orb.Geometry, f ...spatial.Filter) (spr.StandardPlacesResults, error) {

	var poly orb.Geometry

	switch geom.GeoJSONType() {
	case "Polygon":
		poly = geom.(orb.Polygon)
	case "MultiPolygon":
		poly = geom.(orb.MultiPolygon)
	default:
		return nil, fmt.Errorf("Invalid geometry type")
	}

	return db.Intersects(ctx, &poly, f...)
}
