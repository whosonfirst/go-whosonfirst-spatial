package database

import (
	"context"
	"github.com/aaronland/go-roster"
	"github.com/skelterjohn/geom"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-spatial/filter"
	"github.com/whosonfirst/go-spatial/geojson"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"net/url"
)

type SpatialDatabase interface {
	Open(context.Context, string) error
	Close(context.Context) error
	IndexFeature(context.Context, wof_geojson.Feature) error
	GetIntersectsWithCoord(context.Context, geom.Coord, filter.Filter) (spr.StandardPlacesResults, error)
	GetIntersectsWithCoordCandidates(context.Context, geom.Coord) (*geojson.GeoJSONFeatureCollection, error)
	Cache() cache.Cache
}

// PLEASE FIX ME TO HAVE AN OPENER FUNCTION

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

func RegisterSpatialDatabase(ctx context.Context, name string, pr SpatialDatabase) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return indices.Register(ctx, name, pr)
}

func NewSpatialDatabase(ctx context.Context, uri string) (SpatialDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := indices.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	pr := i.(SpatialDatabase)

	err = pr.Open(ctx, uri)

	if err != nil {
		return nil, err
	}

	return pr, nil
}
