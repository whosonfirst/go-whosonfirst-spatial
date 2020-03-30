package database

import (
	"context"
	"github.com/aaronland/go-roster"
	"github.com/skelterjohn/geom"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/geojson"
	"github.com/whosonfirst/go-whosonfirst-spr"
	_ "log"
	"net/url"
)

type SpatialDatabase interface {
	Close(context.Context) error
	IndexFeature(context.Context, wof_geojson.Feature) error
	GetIntersectsWithCoord(context.Context, geom.Coord, filter.Filter) (spr.StandardPlacesResults, error)
	GetIntersectsWithCoordCandidates(context.Context, geom.Coord) (*geojson.GeoJSONFeatureCollection, error)
	ResultsToFeatureCollection(context.Context, spr.StandardPlacesResults) (*geojson.GeoJSONFeatureCollection, error)
}

type SpatialDatabaseInitializeFunc func(ctx context.Context, uri string) (SpatialDatabase, error)

var spatial_databases roster.Roster

func ensureRoster() error {

	if spatial_databases == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		spatial_databases = r
	}

	return nil
}

func RegisterSpatialDatabase(ctx context.Context, scheme string, f SpatialDatabaseInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return spatial_databases.Register(ctx, scheme, f)
}

func NewSpatialDatabase(ctx context.Context, uri string) (SpatialDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := spatial_databases.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(SpatialDatabaseInitializeFunc)
	return f(ctx, uri)
}
