package utils

import (
	"context"
	"encoding/json"
	geojson_utils "github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-spatial/cache"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/geojson"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"io/ioutil"
	_ "log"
)

func IsWOFRecord(fh io.Reader) (bool, error) {

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return false, err
	}

	possible := []string{
		"properties.wof:id",
	}

	id := geojson_utils.Int64Property(body, possible, -1)

	if id == -1 {
		return false, nil
	}

	return true, nil
}

func IsValidRecord(fh io.Reader, ctx context.Context) (bool, error) {

	path, err := index.PathForContext(ctx)

	if err != nil {
		return false, err
	}

	if path == index.STDIN {
		return true, nil
	}

	is_wof, err := uri.IsWOFFile(path)

	if err != nil {
		return false, err
	}

	if !is_wof {
		return false, nil
	}

	is_alt, err := uri.IsAltFile(path)

	if err != nil {
		return false, err
	}

	if is_alt {
		return false, nil
	}

	return true, nil
}

func ResultsToFeatureCollection(ctx context.Context, results spr.StandardPlacesResults, spatial_database database.SpatialDatabase) (*geojson.GeoJSONFeatureCollection, error) {

	c := spatial_database.Cache()

	features := make([]geojson.GeoJSONFeature, 0)

	for _, r := range results.Results() {

		cr, err := c.Get(ctx, r.Id())

		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(cr)

		if err != nil {
			return nil, err
		}

		var fc *cache.FeatureCache

		err = json.Unmarshal(body, &fc)

		if err != nil {
			return nil, err
		}

		f := geojson.GeoJSONFeature{
			Type:       "Feature",
			Properties: fc.SPR(),
			Geometry:   fc.Geometry(),
		}

		features = append(features, f)
	}

	collection := geojson.GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}

	return &collection, nil
}
