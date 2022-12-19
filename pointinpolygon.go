package spatial

import (
	"github.com/paulmach/orb"
)

type PointInPolygonCandidate struct {
	Id        string
	FeatureId string
	IsAlt     bool
	AltLabel  string
	Bounds    orb.Bound
}
