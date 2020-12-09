# go-whosonfirst-spatial

## IMPORTANT

It is work in progress. It works... until it doesn't. It is not well documented yet.

_Once complete this package will supersede the [go-whosonfirst-pip-v2](https://github.com/whosonfirst/go-whosonfirst-pip-v2) package._

## Example

```
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	geojson_utils "github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"	
	_ "github.com/whosonfirst/go-whosonfirst-spatial-rtree"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
)

func main() {

	fl, _ := flags.CommonFlags()
	flags.Parse(fl)

	flags.ValidateCommonFlags(fl)

	paths := fl.Args()
	
	ctx := context.Background()

	spatial_app, _ := app.NewSpatialApplicationWithFlagSet(ctx, fl)
	spatial_app.IndexPaths(ctx, paths...)

	coords, _ := geojson_utils.NewCoordinateFromLatLons(37.794906, -122.395229)
	f, _ := filter.NewSPRFilter()

	spatial_db := spatial_app.SpatialDatabase
	spatial_results, _ := spatial_db.PointInPolygon(ctx, &coords, f)

	body, _ := json.Marshal(spatial_results)
	fmt.Println(string(body))
}
```

_Error handling omitted for brevity._

## Concepts

### Applications

_Please write me_

### Database

_Please write me_

### Filters

_Please write me_

### Indices

_Please write me_

### Standard Places Response (SPR)

_Please write me_

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

* https://github.com/whosonfirst/go-whosonfirst-spatial-rtree
* https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite
* https://github.com/whosonfirst/go-whosonfirst-spatial-http
* https://github.com/whosonfirst/go-whosonfirst-spatial-http-sqlite
* https://github.com/whosonfirst/go-whosonfirst-spatial-grpc