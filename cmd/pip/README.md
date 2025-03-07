# pip

Perform an point-in-polygon operation for an input latitude, longitude coordinate and on a set of Who's on First records stored in a spatial database.

```
$> ./bin/pip -h
Perform an point-in-polygon operation for an input latitude, longitude coordinate and on a set of Who's on First records stored in a spatial database.
Usage:
	 ./bin/pip [options]
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
  -latitude float
    	A valid latitude.
  -longitude float
    	A valid longitude.
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
$> ./bin/pip \
	-latitude 37.617411 \
	-longitude -122.383794 \
	-is-current 1 \
	-iterator-uri repo://#/usr/local/data/sfomuseum-data-maps/ \
	| jq -r '.places[]["wof:name"]'
	
2025/03/07 09:24:56 INFO time to index paths (1) 6.329041ms

SFO (1980)
SFO (1978)
SFO (1965)
SFO (1972)
SFO (1956)
SFO (1960)
SFO (1981)
SFO (1970)
SFO (1988)
SFO (1989)
SFO (1947)
SFO (1998)
```

## Notes

This tool is not especially fast because it is using the default in-memory `rtree://` implementation of the `SpatialDatabase` interface so it will be faster or slower depending on the size and complexity of the source data being indexed.

The "guts" of this application live in the `app/pip` package and are designed such the same application code can be used by other database-specific implementations with a minimal amount of fuss. For example here is the code for the `cmd/pip/main.go` tool in the [whosonfirst/go-whosonfirst-spatial-sqlite](https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite) package:

```
package main

import (
        "context"
        "log"

        _ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
        "github.com/whosonfirst/go-whosonfirst-spatial/app/pip"
)

func main() {

        ctx := context.Background()
        err := pip.Run(ctx)

        if err != nil {
                log.Fatal(err)
        }
}
```