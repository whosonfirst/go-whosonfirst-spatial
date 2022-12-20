package geo

import (
	"testing"
)

func TestNewBoundingBox(t *testing.T) {

	good := [][4]float64{
		[4]float64{-76.0, -60.2, 150.0, 45.0},
	}

	bad := [][4]float64{
		[4]float64{-190.0, -60.2, 150.0, 45.0},
		[4]float64{-170.0, -91.2, 150.0, 45.0},
		[4]float64{-170.0, -51.2, 190.0, 45.0},
		[4]float64{-170.0, -51.2, 120.0, 100.0},
	}

	for _, coords := range good {

		_, err := NewBoundingBox(coords[0], coords[1], coords[2], coords[3])

		if err != nil {
			t.Fatalf("Expected %v to validate, %v", coords, err)
		}
	}

	for _, coords := range bad {

		_, err := NewBoundingBox(coords[0], coords[1], coords[2], coords[3])

		if err == nil {
			t.Fatalf("Expected %v to NOT validate", coords)
		}
	}
}
