package database

import (
	"context"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
)

type ExtrasDatabase interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	Close(context.Context) error
}
