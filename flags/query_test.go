package flags

import (
	"github.com/sfomuseum/go-flags/flagset"
	"testing"
)

func TestAppendQueryFlags(t *testing.T) {

	fs := flagset.NewFlagSet("testing")
	err := AppendQueryFlags(fs)

	if err != nil {
		t.Fatalf("Failed to append query flags, %v", err)
	}
}

func ValidateAppendQueryFlags(t *testing.T) {

	fs := flagset.NewFlagSet("testing")
	err := AppendQueryFlags(fs)

	if err != nil {
		t.Fatalf("Failed to append query flags, %v", err)
	}

	err = ValidateQueryFlags(fs)

	if err != nil {
		t.Fatalf("Failed to validate query flags, %v", err)
	}

}
