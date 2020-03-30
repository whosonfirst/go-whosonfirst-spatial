package database

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spr"
)

type ExtrasDatabase interface {
	IndexFeature(context.Context, geojson.Feature) error
	AppendExtrasWithSPRResults(context.Context, spr.StandardPlacesResults, ...string) error
	Close(context.Context) error
}
