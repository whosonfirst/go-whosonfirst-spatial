package index

import (
	"context"
	"github.com/aaronland/go-roster"
	"github.com/skelterjohn/geom"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-spatial"
	"github.com/whosonfirst/go-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"net/url"
)

type Index interface {
	Open(context.Context, string) error
	Close(context.Context) error
	IndexFeature(context.Context, geojson.Feature) error
	GetIntersectsWithCoord(context.Context, geom.Coord, filter.Filter) (spr.StandardPlacesResults, error)
	GetIntersectsWithCoordCandidates(context.Context, geom.Coord) (*spatial.GeoJSONFeatureCollection, error)
	// GetNearbyWithCoord(geom.Coord, filter.Filter) (spr.StandardPlacesResults, error)
	// GetNearbyWithCoordCandidates(geom.Coord) (*spatial.GeoJSONFeatureCollection, error)
	Cache() cache.Cache
}

type Candidate interface{} // mmmmmaybe?

var indices roster.Roster

func ensureRoster() error {

	if indices == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		indices = r
	}

	return nil
}

func RegisterIndex(ctx context.Context, name string, pr Index) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return indices.Register(ctx, name, pr)
}

func NewIndex(ctx context.Context, uri string) (Index, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := indices.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	pr := i.(Index)

	err = pr.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return pr, nil
}
