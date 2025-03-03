package query

import (
	"context"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type PointInPolygonQuery struct {
	Query
}

func NewPointInPolygonQuery(ctx context.Context, uri string) (Query, error) {
	q := &PointInPolygonQuery{}
	return q, nil
}

func (q *PointInPolygonQuery) Execute(ctx context.Context, db database.SpatialDatabase, geom orb.Geometry, f ...spatial.Filter) (spr.StandardPlacesResults, error) {

	var pt orb.Point

	switch geom.GeoJSONType() {
	case "Point":
		pt = geom.(orb.Point)
	default:
		return nil, fmt.Errorf("Invalid geometry type")
	}

	return db.PointInPolygon(ctx, &pt, f...)
}
