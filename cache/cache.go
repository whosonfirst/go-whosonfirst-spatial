package cache

import (
	"github.com/whosonfirst/go-spatial"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geometry"
	"github.com/whosonfirst/go-whosonfirst-spr"
)

type CacheItem interface {
	SPR() spr.StandardPlacesResult
	Polygons() []geojson.Polygon
	Geometry() spatial.GeoJSONGeometry
}

// see the way we're storing a geojson.Geometry but returning a
// geojson.Polygons per the interface definition? see notes in
// go-whosonfirst-geojson-v2/geometry/polygon.go function called
// PolygonsForFeature for details (20170921/thisisaaronland)

type FeatureCache struct {
	CacheItem       `json:",omitempty"`
	FeatureSPR      spr.StandardPlacesResult `json:"spr"`
	FeaturePolygons []geojson.Polygon        `json:"polygons"`
}

func NewFeatureCache(f geojson.Feature) (CacheItem, error) {

	s, err := f.SPR()

	if err != nil {
		return nil, err
	}

	polys, err := geometry.PolygonsForFeature(f)

	if err != nil {
		return nil, err
	}

	fc := FeatureCache{
		FeatureSPR:      s,
		FeaturePolygons: polys,
	}

	return &fc, nil
}

func (fc *FeatureCache) SPR() spr.StandardPlacesResult {
	return fc.FeatureSPR
}

func (fc *FeatureCache) Geometry() spatial.GeoJSONGeometry {

	multi_poly := make([]spatial.GeoJSONPolygon, 0)

	for _, p := range fc.Polygons() {

		poly := make([]spatial.GeoJSONRing, 0)

		ext := p.ExteriorRing()

		ext_ring := make([]spatial.GeoJSONPoint, 0)

		for _, coord := range ext.Vertices() {

			pt := spatial.GeoJSONPoint{coord.X, coord.Y}
			ext_ring = append(ext_ring, pt)
		}

		poly = append(poly, ext_ring)

		for _, int := range p.InteriorRings() {

			int_ring := make([]spatial.GeoJSONPoint, 0)

			for _, coord := range int.Vertices() {

				pt := spatial.GeoJSONPoint{coord.X, coord.Y}
				int_ring = append(int_ring, pt)
			}

			poly = append(poly, int_ring)
		}

		multi_poly = append(multi_poly, poly)
	}

	geom := spatial.GeoJSONGeometry{
		Type:        "MultiPolygon",
		Coordinates: multi_poly,
	}

	return geom
}

func (fc *FeatureCache) Polygons() []geojson.Polygon {

	return fc.FeaturePolygons
}
