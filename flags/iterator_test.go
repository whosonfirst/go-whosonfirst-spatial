package flags

import (
	"strings"
	"testing"
)

func TestIteratorURIFlag(t *testing.T) {

	str_flag := "repo://#/usr/local/data/sfomuseum-data-architecture"

	fl := new(IteratorURIFlag)

	err := fl.Set(str_flag)

	if err != nil {
		t.Fatalf("Failed to create flag, %v", err)
	}

	if fl.String() != str_flag {
		t.Fatalf("Invalid string value for flag: %s", fl.String())
	}
}

func TestMultiIteratorURIFlag(t *testing.T) {

	str_flags := []string{
		"repo://#/usr/local/data/sfomuseum-data-architecture",
		"featurecollection://#/usr/local/data/sfomuseum.geojson",
	}

	fl := new(MultiIteratorURIFlag)

	for _, str_flag := range str_flags {

		err := fl.Set(str_flag)

		if err != nil {
			t.Fatalf("Failed to add multi flag '%s', %v", str_flag, err)
		}
	}

	if fl.String() != strings.Join(str_flags, SEP_SPACE) {
		t.Fatalf("Invalid string value for flag: %s", fl.String())
	}
}

func TestMultiCSVIteratorURIFlag(t *testing.T) {

	str_flags := []string{
		"repo://#/usr/local/data/sfomuseum-data-architecture",
		"featurecollection://#/usr/local/data/sfomuseum.geojson",
	}

	csv_flag := strings.Join(str_flags, SEP_CSV)

	fl := new(MultiCSVIteratorURIFlag)

	err := fl.Set(csv_flag)

	if err != nil {
		t.Fatalf("Failed to add CSV flag, %v", err)
	}

	if fl.String() != csv_flag {
		t.Fatalf("Invalid string value for flag: %s", fl.String())
	}
}
