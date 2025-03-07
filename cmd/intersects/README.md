# intersects

Perform an intersects operation (as in intersecting geometries) for an input geometry and on a set of Who's on First records stored in a spatial database.

```
$> ./bin/intersects -h
Perform an intersects operation (as in intersecting geometries) for an input geometry and on a set of Who's on First records stored in a spatial database.
Usage:
	 ./bin/intersects [options]
Valid options are:

  -alternate-geometry value
    	One or more alternate geometry labels (wof:alt_label) values to filter results by.
  -cessation string
    	A valid EDTF date string.
  -custom-placetypes string
    	A JSON-encoded string containing custom placetypes defined using the syntax described in the whosonfirst/go-whosonfirst-placetypes repository.
  -enable-custom-placetypes
    	Enable wof:placetype values that are not explicitly defined in the whosonfirst/go-whosonfirst-placetypes repository.
  -geometries string
    	Valid options are: all, alt, default. (default "all")
  -geometry-source string
    	Where to 'read' a geometry (to intersect) from. Valid options are: file, flag, stdin (default "flag")
  -geometry-type string
    	The type of encoding used to perform an intersects operation. Valid options are: geojson, wkt, bbox
  -geometry-value string
    	The value of geometry used to perform an intersects operation. This will vary depending on the value of the -geometry-source flag. For example if -geometry-source=flag then -geometry-value= will be that geometry passed as a string. If -geometry-source=file then -geometry-value= will be the path to a file on disk. If -geometry-source=stdin then -geometry-value will be left empty.
  -inception string
    	A valid EDTF date string.
  -is-ceased value
    	One or more existential flags (-1, 0, 1) to filter results by.
  -is-current value
    	One or more existential flags (-1, 0, 1) to filter results by.
  -is-deprecated value
    	One or more existential flags (-1, 0, 1) to filter results by.
  -is-superseded value
    	One or more existential flags (-1, 0, 1) to filter results by.
  -is-superseding value
    	One or more existential flags (-1, 0, 1) to filter results by.
  -iterator-uri value
    	Zero or more URIs denoting data sources to use for indexing the spatial database at startup. URIs take the form of {ITERATOR_URI} + "#" + {PIPE-SEPARATED LIST OF ITERATOR SOURCES}. Where {ITERATOR_URI} is expected to be a registered whosonfirst/go-whosonfirst-iterate/v2 iterator (emitter) URI and {ITERATOR SOURCES} are valid input paths for that iterator. Supported whosonfirst/go-whosonfirst-iterate/v2 iterator schemes are: cwd://, directory://, featurecollection://, file://, filelist://, geojsonl://, null://, repo://.
  -mode string
    	Valid options are: cli, lambda. (default "cli")
  -placetype value
    	One or more place types to filter results by.
  -properties-reader-uri string
    	A valid whosonfirst/go-reader.Reader URI. Available options are: [fs:// null:// repo:// stdin://]. If the value is {spatial-database-uri} then the value of the '-spatial-database-uri' implements the reader.Reader interface and will be used.
  -property value
    	One or more Who's On First properties to append to each result.
  -sort-uri value
    	Zero or more whosonfirst/go-whosonfirst-spr/sort URIs.
  -spatial-database-uri string
    	A valid whosonfirst/go-whosonfirst-spatial/data.SpatialDatabase URI. options are: [rtree://] (default "rtree://")
  -verbose
    	Enable verbose (debug) logging.
```

## Example

```
$> ./bin/intersects \
	-geometry-source flag \
	-geometry-type bbox \
	-geometry-value '-122.408061,37.601617,-122.354907,37.640167' \
	-iterator-uri repo://#/usr/local/data/sfomuseum-data-whosonfirst \
	| jq -r '.places[]["wof:name"]'
	
2025/03/07 08:55:43 INFO time to index paths (1) 26.923220625s

Earth
94128
Burlingame
San Mateo
San Francisco International Airport
California
United States
North America
```

## Notes

This tool is not especially fast because it is using the default in-memory `rtree://` implementation of the `SpatialDatabase` interface.

The `sfomuseum-data-whosonfirst` has enough records that it's takes a while (26 seconds) to index.

The "guts" of this application live in the `app/intersects` package and are designed such the same application code can be used by other database-specific implementations with a minimal amount of fuss. For example here is the code for the `cmd/intersects/main.go` tool in the [whosonfirst/go-whosonfirst-spatial-sqlite](https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite) package:

```
package main

import (
        "context"
        "log"

        _ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
        "github.com/whosonfirst/go-whosonfirst-spatial/app/intersects"
)

func main() {

        ctx := context.Background()
        err := intersects.Run(ctx)

        if err != nil {
                log.Fatal(err)
        }
}
```