package database

import (
	"context"
	"github.com/aaronland/go-roster"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spatial/geojson"
	"log"
	"net/url"
)

type ExtrasDatabase interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	AppendExtrasWithFeatureCollection(context.Context, *geojson.GeoJSONFeatureCollection, []string) (*geojson.GeoJSONFeatureCollection, error)
	Close(context.Context) error
}

type ExtrasDatabaseInitializeFunc func(ctx context.Context, uri string) (ExtrasDatabase, error)

var extras_databases roster.Roster

func ensureExtrasRoster() error {

	if extras_databases == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		extras_databases = r
	}

	return nil
}

func RegisterExtrasDatabase(ctx context.Context, scheme string, f ExtrasDatabaseInitializeFunc) error {

	err := ensureExtrasRoster()

	if err != nil {
		return err
	}

	log.Println("REGISTER", scheme, f)
	return extras_databases.Register(ctx, scheme, f)
}

func NewExtrasDatabase(ctx context.Context, uri string) (ExtrasDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := extras_databases.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(ExtrasDatabaseInitializeFunc)
	return f(ctx, uri)
}
