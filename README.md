# go-whosonfirst-spatial

![](docs/images/wof-spatial-sfo.png)

## IMPORTANT

It is work in progress. It works... until it doesn't. It is not documented.

_Once complete this package will supersede the [go-whosonfirst-pip-v2](https://github.com/whosonfirst/go-whosonfirst-pip-v2) package._

## Interfaces

### SpatialDatabase

```
type SpatialDatabase interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	PointInPolygon(context.Context, geom.Coord, filter.Filter) (spr.StandardPlacesResults, error)
	PointInPolygonCandidates(context.Context, geom.Coord) (*geojson.GeoJSONFeatureCollection, error)
	PointInPolygonWithChannels(context.Context, geom.Coord, filter.Filter, chan spr.StandardPlacesResult, chan error, chan bool)
	PointInPolygonCandidatesWithChannels(context.Context, geom.Coord, chan geojson.GeoJSONFeature, chan error, chan bool)
	StandardPlacesResultsToFeatureCollection(context.Context, spr.StandardPlacesResults) (*geojson.GeoJSONFeatureCollection, error)
	Close(context.Context) error
}
```

### PropertiesReader

```
type PropertiesReader interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	PropertiesResponseResultsWithStandardPlacesResults(context.Context, spr.StandardPlacesResults, []string) (*PropertiesResponseResults, error)
	AppendPropertiesWithFeatureCollection(context.Context, *geojson.GeoJSONFeatureCollection, []string) error
	Close(context.Context) error
}
```

## See also