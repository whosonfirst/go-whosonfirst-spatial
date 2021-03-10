package api

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"net/url"
	"strconv"
)

type PointInPolygonRequest struct {
	Latitude            float64  `json:"latitude"`
	Longitude           float64  `json:"longitude"`
	Properties          []string `json:"properties"`
	Placetypes          []string `json:"placetypes,omitempty"`
	Geometries          string   `json:"geometries,omitempty"`
	AlternateGeometries []string `json:"alternate_geometries,omitempty"`
	IsCurrent           []int    `json:"is_current,omitempty"`
	IsCeased            []int    `json:"is_ceased,omitempty"`
	IsDeprecated        []int    `json:"is_deprecated,omitempty"`
	IsSuperseded        []int    `json:"is_superseded,omitempty"`
	IsSuperseding       []int    `json:"is_superseding,omitempty"`
}

func NewPointInPolygonRequestFromFlagSet(fs *flag.FlagSet) (*PointInPolygonRequest, error) {

	req := &PointInPolygonRequest{}

	latitude, err := flags.Float64Var(fs, "latitude")

	if err != nil {
		return nil, err
	}

	req.Latitude = latitude

	longitude, err := flags.Float64Var(fs, "longitude")

	if err != nil {
		return nil, err
	}

	req.Longitude = longitude

	props, err := flags.MultiStringVar(fs, "properties")

	if err != nil {
		return nil, err
	}

	req.Properties = props

	placetypes, err := flags.MultiStringVar(fs, "placetype")

	if err != nil {
		return nil, err
	}

	req.Placetypes = placetypes

	geometries, err := flags.StringVar(fs, "geometries")

	if err != nil {
		return nil, err
	}

	req.Geometries = geometries

	alt_geoms, err := flags.MultiStringVar(fs, "alternate-geometry")

	if err != nil {
		return nil, err
	}

	req.AlternateGeometries = alt_geoms

	is_current, err := flags.MultiIntVar(fs, "is-current")

	if err != nil {
		return nil, err
	}

	req.IsCurrent = is_current

	is_ceased, err := flags.MultiIntVar(fs, "is-ceased")

	if err != nil {
		return nil, err
	}

	req.IsCeased = is_ceased

	is_deprecated, err := flags.MultiIntVar(fs, "is-deprecated")

	if err != nil {
		return nil, err
	}

	req.IsDeprecated = is_deprecated

	is_superseded, err := flags.MultiIntVar(fs, "is-superseded")

	if err != nil {
		return nil, err
	}

	req.IsSuperseded = is_superseded

	is_superseding, err := flags.MultiIntVar(fs, "is-superseding")

	if err != nil {
		return nil, err
	}

	req.IsSuperseding = is_superseding

	return req, nil
}

func NewSPRFilterFromPointInPolygonRequest(req *PointInPolygonRequest) (filter.Filter, error) {

	q := url.Values{}
	q.Set("geometries", req.Geometries)

	for _, v := range req.AlternateGeometries {
		q.Add("alternate_geometry", v)
	}

	for _, v := range req.Placetypes {
		q.Add("placetype", v)
	}

	for _, v := range req.IsCurrent {
		q.Add("is_current", strconv.Itoa(v))
	}

	for _, v := range req.IsCeased {
		q.Add("is_ceased", strconv.Itoa(v))
	}

	for _, v := range req.IsDeprecated {
		q.Add("is_deprecated", strconv.Itoa(v))
	}

	for _, v := range req.IsSuperseded {
		q.Add("is_superseded", strconv.Itoa(v))
	}

	for _, v := range req.IsSuperseding {
		q.Add("is_superseding", strconv.Itoa(v))
	}

	return filter.NewSPRFilterFromQuery(q)
}
