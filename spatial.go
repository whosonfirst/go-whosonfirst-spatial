package spatial

import (
	"context"
	"github.com/skelterjohn/geom"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spr"
)

type SpatialIndex interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	PointInPolygon(context.Context, *geom.Coord, ...filter.Filter) (spr.StandardPlacesResults, error)
	PointInPolygonCandidates(context.Context, *geom.Coord, ...filter.Filter) ([]*PointInPolygonCandidate, error)
	PointInPolygonWithChannels(context.Context, chan spr.StandardPlacesResult, chan error, chan bool, *geom.Coord, ...filter.Filter)
	PointInPolygonCandidatesWithChannels(context.Context, chan *PointInPolygonCandidate, chan error, chan bool, *geom.Coord, ...filter.Filter)
	Disconnect(context.Context) error
}

type PointInPolygonCandidate struct {
	Id        string
	FeatureId string
	IsAlt     bool
	AltLabel  string
	Bounds    *geom.Rect
}

type PropertiesResponse map[string]interface{}

type PropertiesResponseResults struct {
	Properties []*PropertiesResponse `json:"properties"`
}
