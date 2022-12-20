package geo

import (
	"testing"
)

func TestIsValidLatitude(t *testing.T) {

	good := []float64{
		-45.0,
		67.3,
	}

	bad := []float64{
		-145.0,
		167.3,
	}

	for _, lat := range good {

		ok := IsValidLatitude(lat)

		if !ok {
			t.Fatalf("Expected %f to validate", lat)
		}
	}

	for _, lat := range bad {

		ok := IsValidLatitude(lat)

		if ok {
			t.Fatalf("Expected %f to NOT validate", lat)
		}
	}

}

func TestIsValidLongitude(t *testing.T) {

	good := []float64{
		-145.0,
		167.3,
	}

	bad := []float64{
		-245.0,
		197.3,
	}

	for _, lon := range good {

		ok := IsValidLongitude(lon)

		if !ok {
			t.Fatalf("Expected %f to validate", lon)
		}
	}

	for _, lon := range bad {

		ok := IsValidLongitude(lon)

		if ok {
			t.Fatalf("Expected %f to NOT validate", lon)
		}
	}

}

func TestIsValidMinLatitude(t *testing.T) {

	good := []float64{
		-45.0,
	}

	bad := []float64{
		-145.0,
	}

	for _, lat := range good {

		ok := IsValidMinLatitude(lat)

		if !ok {
			t.Fatalf("Expected %f to validate", lat)
		}
	}

	for _, lat := range bad {

		ok := IsValidMinLatitude(lat)

		if ok {
			t.Fatalf("Expected %f to NOT validate", lat)
		}
	}

}

func TestIsValidMaxLatitude(t *testing.T) {

	good := []float64{
		67.3,
	}

	bad := []float64{
		167.3,
	}

	for _, lat := range good {

		ok := IsValidMaxLatitude(lat)

		if !ok {
			t.Fatalf("Expected %f to validate", lat)
		}
	}

	for _, lat := range bad {

		ok := IsValidMaxLatitude(lat)

		if ok {
			t.Fatalf("Expected %f to NOT validate", lat)
		}
	}

}

func TestIsValidMinLongitude(t *testing.T) {

	good := []float64{
		167.3,
	}

	bad := []float64{
		197.3,
	}

	for _, lon := range good {

		ok := IsValidMinLongitude(lon)

		if !ok {
			t.Fatalf("Expected %f to validate", lon)
		}
	}

	for _, lon := range bad {

		ok := IsValidMinLongitude(lon)

		if ok {
			t.Fatalf("Expected %f to NOT validate", lon)
		}
	}

}

func TestIsValidMaxLongitude(t *testing.T) {

	good := []float64{
		-145.0,
	}

	bad := []float64{
		-245.0,
	}

	for _, lon := range good {

		ok := IsValidMaxLongitude(lon)

		if !ok {
			t.Fatalf("Expected %f to validate", lon)
		}
	}

	for _, lon := range bad {

		ok := IsValidMaxLongitude(lon)

		if ok {
			t.Fatalf("Expected %f to NOT validate", lon)
		}
	}

}
