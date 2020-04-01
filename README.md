# go-whosonfirst-spatial

![](docs/images/wof-spatial-sfo.png)

## IMPORTANT

DO NOT TRY TO USE THIS YET. It is work in progress.

_Once complete this package will supersede the [go-whosonfirst-pip-v2](https://github.com/whosonfirst/go-whosonfirst-pip-v2) package._

## Interfaces

### SpatialDatabase

```
type SpatialDatabase interface {
	Close(context.Context) error
	IndexFeature(context.Context, wof_geojson.Feature) error
	PointInPolygon(context.Context, geom.Coord, filter.Filter) (spr.StandardPlacesResults, error)
	PointInPolygonCandidates(context.Context, geom.Coord) (*geojson.GeoJSONFeatureCollection, error)
	PointInPolygonWithChannels(context.Context, geom.Coord, filter.Filter, chan spr.StandardPlacesResult, chan error, chan bool)
	PointInPolygonCandidatesWithChannels(context.Context, geom.Coord, chan geojson.GeoJSONFeature, chan error, chan bool)
	ResultsToFeatureCollection(context.Context, spr.StandardPlacesResults) (*geojson.GeoJSONFeatureCollection, error)
}
```

### ExtrasReader

```
type ExtrasReader interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	AppendExtras(context.Context, interface{}, []string) error
	AppendExtrasWithStandardPlacesResults(context.Context, spr.StandardPlacesResults, []string) error
	AppendExtrasWithFeatureCollection(context.Context, *geojson.GeoJSONFeatureCollection, []string) error
	Close(context.Context) error
}
```

