package utils

import (
	"context"
	"encoding/json"
	"github.com/whosonfirst/go-spatial"
	"github.com/whosonfirst/go-spatial/cache"
	pip_index "github.com/whosonfirst/go-spatial/index"
	geojson_utils "github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"github.com/whosonfirst/go-whosonfirst-index"
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

func ResultsToFeatureCollection(ctx context.Context, results spr.StandardPlacesResults, idx pip_index.Index) (*spatial.GeoJSONFeatureCollection, error) {

	c := idx.Cache()

	features := make([]spatial.GeoJSONFeature, 0)

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

		f := spatial.GeoJSONFeature{
			Type:       "Feature",
			Properties: fc.SPR(),
			Geometry:   fc.Geometry(),
		}

		features = append(features, f)
	}

	collection := spatial.GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}

	return &collection, nil
}
