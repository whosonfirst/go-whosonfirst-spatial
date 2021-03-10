package api

import (
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
