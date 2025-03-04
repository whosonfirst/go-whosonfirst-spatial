package spatial

import (
	"context"
	"iter"

	"github.com/paulmach/orb"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

// SpatialIndex defines an interface for performing spatial operations on a collection of Who's On First GeoJSON features.
type SpatialIndex interface {
	// IndexFeature adds a Who's On First GeoJSON feature to the index.
	IndexFeature(context.Context, []byte) error
	// RemoveFeature removes a Who's On First GeoJSON feature from the index.
	RemoveFeature(context.Context, string) error
	// PointInPolygon performs a point-in-polygon operation to retrieve matching records from the index.
	PointInPolygon(context.Context, *orb.Point, ...Filter) (spr.StandardPlacesResults, error)
	PointInPolygonWithIterator(context.Context, *orb.Point, ...Filter) iter.Seq2[spr.StandardPlacesResult, error]
	Intersects(context.Context, orb.Geometry, ...Filter) (spr.StandardPlacesResults, error)
	IntersectsWithIterator(context.Context, orb.Geometry, ...Filter) iter.Seq2[spr.StandardPlacesResult, error]
	// Disconnect closes any underlying connections used by the index.
	Disconnect(context.Context) error
}
