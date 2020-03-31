package api

import (
	"encoding/json"
	"github.com/aaronland/go-http-sanitize"
	geojson_utils "github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	_ "log"
	"net/http"
	"strconv"
	"strings"
)

type PointInPolygonHandlerOptions struct {
	EnableGeoJSON bool
}

func PointInPolygonHandler(spatial_app *app.SpatialApplication, opts *PointInPolygonHandlerOptions) (http.Handler, error) {

	spatial_db := spatial_app.SpatialDatabase
	extras_db := spatial_app.ExtrasDatabase
	walker := spatial_app.Walker

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		if walker.IsIndexing() {
			http.Error(rsp, "indexing records", http.StatusServiceUnavailable)
			return
		}

		ctx := req.Context()
		query := req.URL.Query()

		str_lat, err := sanitize.GetString(req, "latitude")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		str_lon, err := sanitize.GetString(req, "longitude")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		str_format, err := sanitize.GetString(req, "format")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if str_format == "geojson" && !opts.EnableGeoJSON {
			http.Error(rsp, "Invalid format", http.StatusBadRequest)
			return
		}

		if str_lat == "" {
			http.Error(rsp, "Missing 'latitude' parameter", http.StatusBadRequest)
			return
		}

		if str_lon == "" {
			http.Error(rsp, "Missing 'longitude' parameter", http.StatusBadRequest)
			return
		}

		lat, err := strconv.ParseFloat(str_lat, 64)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		lon, err := strconv.ParseFloat(str_lon, 64)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		coord, err := geojson_utils.NewCoordinateFromLatLons(lat, lon)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		filters, err := filter.NewSPRFilterFromQuery(query)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		results, err := spatial_db.PointInPolygon(ctx, coord, filters)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		var final interface{}
		final = results

		if str_format == "geojson" {

			collection, err := spatial_db.ResultsToFeatureCollection(ctx, results)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			final = collection
		}

		js, err := json.Marshal(final)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
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
					http.Error(rsp, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(js)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}

func PointInPolygonCandidatesHandler(spatial_app *app.SpatialApplication) (http.Handler, error) {

	walker := spatial_app.Walker
	spatial_db := spatial_app.SpatialDatabase

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		if walker.IsIndexing() {
			http.Error(rsp, "indexing records", http.StatusServiceUnavailable)
			return
		}

		ctx := req.Context()

		str_lat, err := sanitize.GetString(req, "latitude")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		str_lon, err := sanitize.GetString(req, "longitude")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if str_lat == "" {
			http.Error(rsp, "Missing 'latitude' parameter", http.StatusBadRequest)
			return
		}

		if str_lon == "" {
			http.Error(rsp, "Missing 'longitude' parameter", http.StatusBadRequest)
			return
		}

		lat, err := strconv.ParseFloat(str_lat, 64)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		lon, err := strconv.ParseFloat(str_lon, 64)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		coord, err := geojson_utils.NewCoordinateFromLatLons(lat, lon)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		candidates, err := spatial_db.PointInPolygonCandidates(ctx, coord)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		enc, err := json.Marshal(candidates)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(enc)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
