package api

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
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

func (req *PointInPolygonRequest) AppendSPRFilterFromFlagSet(fs *flag.FlagSet) error {

	placetypes, err := MultiStringVar(fs, "placetype")

	if err != nil {
		return err
	}

	req.Placetypes = placetypes

	geometries, err := StringVar(fs, "geometries")

	if err != nil {
		return err
	}

	req.Geometries = geometries

	alt_geoms, err := MultiStringVar(fs, "alternate-geometry")

	if err != nil {
		return err
	}

	req.AlternateGeometries = alt_geom

	props, err := MultiStringVar(fs, "properties")

	if err != nil {
		return err
	}

	req.Properties = props

	is_current, err := MultiStringVar(fs, "is-current")

	if err != nil {
		return err
	}

	is_ceased, err := MultiStringVar(fs, "is-ceased")

	if err != nil {
		return err
	}

	is_deprecated, err := MultiStringVar(fs, "is-deprecated")

	if err != nil {
		return err
	}

	is_superseded, err := MultiStringVar(fs, "is-superseded")

	if err != nil {
		return err
	}

	is_superseding, err := MultiStringVar(fs, "is-superseding")

	if err != nil {
		return err
	}

	return nil
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
