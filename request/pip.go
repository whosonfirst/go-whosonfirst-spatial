package request

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

func (q *PointInPolygonQuery) Execute(ctx context.Context, db database.SpatialDatabase, geom orb.Geometry, f  spatial.Filter) (spr.StandardPlacesResults, error) {

	return nil, fmt.Errorf("Not implemented")
}
