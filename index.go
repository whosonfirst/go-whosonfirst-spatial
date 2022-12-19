package spatial

import (
	"context"
	"github.com/paulmach/orb"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type SpatialIndex interface {
	IndexFeature(context.Context, []byte) error
	RemoveFeature(context.Context, string) error
	PointInPolygon(context.Context, *orb.Point, ...Filter) (spr.StandardPlacesResults, error)
	PointInPolygonCandidates(context.Context, *orb.Point, ...Filter) ([]*PointInPolygonCandidate, error)
	PointInPolygonWithChannels(context.Context, chan spr.StandardPlacesResult, chan error, chan bool, *orb.Point, ...Filter)
	PointInPolygonCandidatesWithChannels(context.Context, chan *PointInPolygonCandidate, chan error, chan bool, *orb.Point, ...Filter)
	Disconnect(context.Context) error
}
