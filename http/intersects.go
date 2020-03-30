package http

import (
	"encoding/json"
	"errors"
	geojson_utils "github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"html/template"
	_ "log"
	gohttp "net/http"
	"strconv"
	"strings"
)

type IntersectsHandlerOptions struct {
	EnableGeoJSON bool
	Templates     *template.Template
}

type IntersectsWWWHandlerOptions struct {
	Templates *template.Template
}

func IntersectsWWWHandler(spatial_app *app.SpatialApplication, opts *IntersectsWWWHandlerOptions) (gohttp.Handler, error) {

	t := opts.Templates.Lookup("intersects")

	if t == nil {
		return nil, errors.New("Missing intersects template")
	}

	walker := spatial_app.Walker

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		if walker.IsIndexing() {
			gohttp.Error(rsp, "indexing records", gohttp.StatusServiceUnavailable)
			return
		}

		// important if we're trying to use this in a Lambda/API Gateway context

		rsp.Header().Set("Content-Type", "text/html; charset=utf-8")

		err := t.Execute(rsp, nil)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}

func IntersectsHandler(spatial_app *app.SpatialApplication, opts *IntersectsHandlerOptions) (gohttp.Handler, error) {

	spatial_db := spatial_app.SpatialDatabase
	extras_db := spatial_app.ExtrasDatabase
	walker := spatial_app.Walker

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		if walker.IsIndexing() {
			gohttp.Error(rsp, "indexing records", gohttp.StatusServiceUnavailable)
			return
		}

		ctx := req.Context()
		query := req.URL.Query()

		str_lat := query.Get("latitude")
		str_lon := query.Get("longitude")
		str_format := query.Get("format")

		if str_format == "geojson" && !opts.EnableGeoJSON {
			gohttp.Error(rsp, "Invalid format", gohttp.StatusBadRequest)
			return
		}

		if str_lat == "" {
			gohttp.Error(rsp, "Missing 'latitude' parameter", gohttp.StatusBadRequest)
			return
		}

		if str_lon == "" {
			gohttp.Error(rsp, "Missing 'longitude' parameter", gohttp.StatusBadRequest)
			return
		}

		lat, err := strconv.ParseFloat(str_lat, 64)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		lon, err := strconv.ParseFloat(str_lon, 64)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		coord, err := geojson_utils.NewCoordinateFromLatLons(lat, lon)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		filters, err := filter.NewSPRFilterFromQuery(query)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		results, err := spatial_db.GetIntersectsWithCoord(ctx, coord, filters)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		var final interface{}
		final = results

		if str_format == "geojson" {

			collection, err := spatial_db.ResultsToFeatureCollection(ctx, results)

			if err != nil {
				gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
				return
			}

			final = collection
		}

		js, err := json.Marshal(final)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		// experimental - see notes in extras/extras.go (20180303/thisisaaronland)

		if extras_db != nil {

			var extras_paths []string

			str_extras := query.Get("extras")
			str_extras = strings.Trim(str_extras, " ")

			if str_extras != "" {
				extras_paths = strings.Split(str_extras, ",")
			}

			if len(extras_paths) > 0 {

				// FIX ME - WRITE TO JS OR ... ?

				err := extras_db.AppendExtrasWithSPRResults(ctx, results, extras_paths...)

				if err != nil {
					gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
					return
				}
			}
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(js)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}

func IntersectsCandidatesHandler(spatial_app *app.SpatialApplication) (gohttp.Handler, error) {

	walker := spatial_app.Walker
	spatial_db := spatial_app.SpatialDatabase

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		if walker.IsIndexing() {
			gohttp.Error(rsp, "indexing records", gohttp.StatusServiceUnavailable)
			return
		}

		ctx := req.Context()
		query := req.URL.Query()

		str_lat := query.Get("latitude")
		str_lon := query.Get("longitude")

		if str_lat == "" {
			gohttp.Error(rsp, "Missing 'latitude' parameter", gohttp.StatusBadRequest)
			return
		}

		if str_lon == "" {
			gohttp.Error(rsp, "Missing 'longitude' parameter", gohttp.StatusBadRequest)
			return
		}

		lat, err := strconv.ParseFloat(str_lat, 64)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		lon, err := strconv.ParseFloat(str_lon, 64)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		coord, err := geojson_utils.NewCoordinateFromLatLons(lat, lon)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		candidates, err := spatial_db.GetIntersectsWithCoordCandidates(ctx, coord)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		enc, err := json.Marshal(candidates)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(enc)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
