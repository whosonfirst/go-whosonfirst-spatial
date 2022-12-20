package geo

import (
	"fmt"
	"github.com/paulmach/orb"
)

// NewCoordinate returns a new `orb.Point` instance derived from 'x' and 'y'.
func NewCoordinate(x float64, y float64) (*orb.Point, error) {

	if !IsValidLongitude(x) {
		return nil, fmt.Errorf("Invalid longitude")
	}

	if !IsValidLatitude(y) {
		return nil, fmt.Errorf("Invalid latitude")
	}

	coord := &orb.Point{x, y}
	return coord, nil
}
