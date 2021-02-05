package geo

import (
	"github.com/paulmach/orb"
)

func MultiPolygonContainsCoord(multi [][][][]float64, c *orb.Point) bool {

	for _, poly := range multi {

		if PolygonContainsCoord(poly, c) {
			return true
		}
	}

	return false
}

func PolygonContainsCoord(poly [][][]float64, c *orb.Point) bool {

	count := len(poly)

	if count == 0 {
		return false
	}

	// exterior ring

	exterior_ring := poly[0]

	if !RingContainsCoord(exterior_ring, c) {
		return false
	}

	// interior rings

	if count > 1 {

		for _, interior_ring := range poly {

			if RingContainsCoord(interior_ring, c) {
				return false
			}
		}
	}

	return true
}

// FIX ME...

func RingContainsCoord(ring [][]float64, c *orb.Point) bool {

	polygon := geom.Polygon{}

	for _, pt := range ring {
		polygon.AddVertex(geom.Coord{X: pt[0], Y: pt[1]})
	}

	return polygon.ContainsCoord(*c)
}
