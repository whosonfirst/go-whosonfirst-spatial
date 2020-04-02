package extras

import (
	"context"
	"github.com/aaronland/go-roster"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spatial/geojson"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"net/url"
)

type Properties map[string]interface{}

type PropertiesResponse struct {
	Properties []*Properties `json:"properties"`
}

type ExtrasReader interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	PropertiesResponseWithStandardPlacesResults(context.Context, spr.StandardPlacesResults, []string) (*PropertiesResponse, error)
	AppendPropertiesWithFeatureCollection(context.Context, *geojson.GeoJSONFeatureCollection, []string) error
	Close(context.Context) error
}

type ExtrasReaderInitializeFunc func(ctx context.Context, uri string) (ExtrasReader, error)

var extras_readers roster.Roster

func ensureExtrasRoster() error {

	if extras_readers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		extras_readers = r
	}

	return nil
}

func RegisterExtrasReader(ctx context.Context, scheme string, f ExtrasReaderInitializeFunc) error {

	err := ensureExtrasRoster()

	if err != nil {
		return err
	}

	return extras_readers.Register(ctx, scheme, f)
}

func NewExtrasReader(ctx context.Context, uri string) (ExtrasReader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := extras_readers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(ExtrasReaderInitializeFunc)
	return f(ctx, uri)
}
