package cache

import (
	"github.com/paulmach/go.geojson"
	"github.com/whosonfirst/go-whosonfirst-spr"
)

type CacheItem interface {
	SPR() (spr.StandardPlacesResult, error)
	Geometry() (*geojson.Geometry, error)
}
