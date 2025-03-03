package request

import (
	"context"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/whosonfirst/go-whosonfirst-spatial"	
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	// "github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-spr/v2/sort"
)

// const timingsPIPQuery string = "PIP query"
// const timingsPIPQueryPointInPolygon string = "PIP query point in polygon"
// const timingsPIPQuerySort string = "PIP query sort"

type Query interface {
	Execute(context.Context, database.SpatialDatabase, orb.Geometry, spatial.Filter) (spr.StandardPlacesResults, error) 
}

func ExecuteQuery(ctx context.Context, db database.SpatialDatabase, q Query, geom orb.Geometry, req *SpatialRequest) (spr.StandardPlacesResults, error) {

	f, err := NewSPRFilterFromSpatialRequest(req)

	if err != nil {
		return nil, fmt.Errorf("Failed to create point in polygon filter from request, %w", err)
	}

	var principal_sorter sort.Sorter
	var follow_on_sorters []sort.Sorter

	for idx, uri := range req.Sort {

		s, err := sort.NewSorter(ctx, uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create sorter for '%s', %w", uri, err)
		}

		if idx == 0 {
			principal_sorter = s
		} else {
			follow_on_sorters = append(follow_on_sorters, s)
		}
	}

	rsp, err := q.Execute(ctx, db, geom, f)

	if err != nil {
		return nil, fmt.Errorf("Failed to perform point in polygon query, %w", err)
	}

	if principal_sorter != nil {

		sorted, err := principal_sorter.Sort(ctx, rsp, follow_on_sorters...)

		if err != nil {
			return nil, fmt.Errorf("Failed to sort results, %w", err)
		}

		rsp = sorted
	}

	return rsp, nil
}
