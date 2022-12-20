package geo

import (
	"fmt"
	"github.com/paulmach/orb"
)

// NewBoundingBox returns a new `orb.Bound` instance derived from 'minx', 'miny', 'maxx' and 'maxy'
func NewBoundingBox(minx float64, miny float64, maxx float64, maxy float64) (*orb.Bound, error) {

	if !IsValidLongitude(minx) {
		return nil, fmt.Errorf("Invalid min longitude")
	}

	if !IsValidLatitude(miny) {
		return nil, fmt.Errorf("Invalid min latitude")
	}

	if !IsValidLongitude(maxx) {
		return nil, fmt.Errorf("Invalid max longitude")
	}

	if !IsValidLatitude(maxy) {
		return nil, fmt.Errorf("Invalid max latitude")
	}

	if minx > maxx {
		return nil, fmt.Errorf("Min lon is greater than max lon")
	}

	if minx > maxx {
		return nil, fmt.Errorf("Min latitude is greater than max latitude")
	}

	min_coord, err := NewCoordinate(minx, miny)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new (SW) coordinate, %w", err)
	}

	max_coord, err := NewCoordinate(maxx, maxy)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new (NE) coordinate, %w", err)
	}

	rect := &orb.Bound{
		Min: *min_coord,
		Max: *max_coord,
	}

	return rect, nil
}
