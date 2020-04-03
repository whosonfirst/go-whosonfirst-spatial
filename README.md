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

## Tools

### spatial-server

```
go run -mod vendor cmd/spatial-server/main.go \
	-enable-www \
	-spatial-database-uri 'rtree:///?strict=false' \
	-properties-reader-uri 'whosonfirst:///?reader=fs:///usr/local/data/sfomuseum-data-maps/data&cache=gocache://' \
	-data-endpoint https://millsfield.sfomuseum.org/data \
	-nextzen-apikey {APIKEY}
	-mode directory:// /usr/local/data/sfomuseum-data-maps/data
	
2020/04/03 09:40:26 -enable-www flag is true causing the following flags to also be true: -enable-geojson -enable-candidates -enable-properties
2020/04/03 09:40:26 Feature ID 1360391313 triggered the following warning: Invalid wof:placetype
2020/04/03 09:40:26 Feature ID 1360391317 triggered the following warning: Invalid wof:placetype
2020/04/03 09:40:26 Feature ID 1360391321 triggered the following warning: Invalid wof:placetype
...
09:40:26.241237 [main] STATUS finished indexing in 27.925694ms
```

## See also