package spatial

import (
	"github.com/skelterjohn/geom"
)

type PointInPolygonCandidate struct {
	Id int64
	WOFId string
	AltLabel string
	Bounds *geom.Rect
}
