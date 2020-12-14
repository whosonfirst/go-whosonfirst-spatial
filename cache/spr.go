package cache

import (
	"github.com/paulmach/go.geojson"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/geometry"
	"github.com/whosonfirst/go-whosonfirst-spr"
)

// see the way we're storing a geojson.Geometry but returning a
// geojson.Polygons per the interface definition? see notes in
// go-whosonfirst-geojson-v2/geometry/polygon.go function called
// PolygonsForFeature for details (20170921/thisisaaronland)

type SPRCacheItem struct {
	CacheItem `json:",omitempty"`
	feature   wof_geojson.Feature
}

func NewSPRCacheItem(f wof_geojson.Feature) (CacheItem, error) {

	fc := SPRCacheItem{
		feature: f,
	}

	return &fc, nil
}

func (fc *SPRCacheItem) SPR() (spr.StandardPlacesResult, error) {
	return fc.feature.SPR()
}

func (fc *SPRCacheItem) Geometry() (*geojson.Geometry, error) {

	str_geom, err := geometry.ToString(fc.feature)

	if err != nil {
		return nil, err
	}

	return geojson.UnmarshalGeometry([]byte(str_geom))
}
