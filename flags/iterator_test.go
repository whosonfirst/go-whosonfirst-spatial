package flags

import (
	"github.com/sfomuseum/go-flags/flagset"
	"testing"
)

func TestAppendIndexingFlags(t *testing.T) {

	fs := flagset.NewFlagSet("testing")
	err := AppendIndexingFlags(fs)

	if err != nil {
		t.Fatalf("Failed to append indexing flags, %v", err)
	}
}

func ValidateAppendIndexingFlags(t *testing.T) {

	fs := flagset.NewFlagSet("testing")
	err := AppendIndexingFlags(fs)

	if err != nil {
		t.Fatalf("Failed to append indexing flags, %v", err)
	}

	err = ValidateIndexingFlags(fs)

	if err != nil {
		t.Fatalf("Failed to validate indexing flags, %v", err)
	}

}
