package geo

import (
	"testing"
)

func TestNewCoordinate(t *testing.T) {

	good := [][2]float64{
		[2]float64{-76.0, -60.2},
	}

	bad := [][2]float64{
		[2]float64{-190.0, -60.2},
		[2]float64{-170.0, -91.2},
	}

	for _, coords := range good {

		_, err := NewCoordinate(coords[0], coords[1])

		if err != nil {
			t.Fatalf("Expected %v to validate, %v", coords, err)
		}
	}

	for _, coords := range bad {

		_, err := NewCoordinate(coords[0], coords[1])

		if err == nil {
			t.Fatalf("Expected %v to NOT validate", coords)
		}
	}
}
