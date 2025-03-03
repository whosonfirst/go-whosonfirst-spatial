package request

import (

)

type SpatialRequest struct {
	Placetypes          []string `json:"placetypes,omitempty"`
	Geometries          string   `json:"geometries,omitempty"`
	AlternateGeometries []string `json:"alternate_geometries,omitempty"`
	IsCurrent           []int64  `json:"is_current,omitempty"`
	IsCeased            []int64  `json:"is_ceased,omitempty"`
	IsDeprecated        []int64  `json:"is_deprecated,omitempty"`
	IsSuperseded        []int64  `json:"is_superseded,omitempty"`
	IsSuperseding       []int64  `json:"is_superseding,omitempty"`
	InceptionDate       string   `json:"inception_date,omitempty"`
	CessationDate       string   `json:"cessation_date,omitempty"`
	Properties          []string `json:"properties,omitempty"`
	Sort                []string `json:"sort,omitempty"`
}
