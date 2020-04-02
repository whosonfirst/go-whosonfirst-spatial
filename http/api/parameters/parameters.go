package parameters

import (
	"errors"
	"github.com/aaronland/go-http-sanitize"
	"github.com/skelterjohn/geom"
	geojson_utils "github.com/whosonfirst/go-whosonfirst-geojson-v2/utils"
	"net/http"
	"strconv"
	"strings"
)

func Latitude(req *http.Request) (float64, error) {

	str_lat, err := sanitize.GetString(req, "latitude")

	if err != nil {
		return 0, err
	}

	if str_lat == "" {
		return 0, errors.New("Missing 'latitude' parameter")
	}

	lat, err := strconv.ParseFloat(str_lat, 64)

	if err != nil {
		return 0, err
	}

	if lat > 90.0 || lat < -90.0 {
		return 0, errors.New("Invalid 'latitude' parameter")
	}

	return lat, nil
}

func Longitude(req *http.Request) (float64, error) {

	str_lat, err := sanitize.GetString(req, "longitude")

	if err != nil {
		return 0, err
	}

	if str_lat == "" {
		return 0, errors.New("Missing 'longitude' parameter")
	}

	lat, err := strconv.ParseFloat(str_lat, 64)

	if err != nil {
		return 0, err
	}

	if lat > 180.0 || lat < -180.0 {
		return 0, errors.New("Invalid 'longitude' parameter")
	}

	return lat, nil
}

func Coordinate(req *http.Request) (*geom.Coord, error) {

	lat, err := Latitude(req)

	if err != nil {
		return nil, err
	}

	lon, err := Longitude(req)

	if err != nil {
		return nil, err
	}

	coord, err := geojson_utils.NewCoordinateFromLatLons(lat, lon)

	if err != nil {
		return nil, err
	}

	return &coord, err
}

func Properties(req *http.Request) ([]string, error) {

	var properties []string

	str_properties, err := sanitize.GetString(req, "properties")

	if err != nil {
		return nil, err
	}

	str_properties = strings.Trim(str_properties, " ")

	if str_properties != "" {
		properties = strings.Split(str_properties, ",")
	}

	return properties, nil
}
