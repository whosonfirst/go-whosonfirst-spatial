package flags

import (
	"github.com/sfomuseum/go-flags/flagset"
	"testing"
)

func TestAppendCommonFlags(t *testing.T) {

	fs := flagset.NewFlagSet("testing")
	err := AppendCommonFlags(fs)

	if err != nil {
		t.Fatalf("Failed to append common flags, %v", err)
	}
}

func ValidateAppendCommonFlags(t *testing.T) {

	fs := flagset.NewFlagSet("testing")
	err := AppendCommonFlags(fs)

	if err != nil {
		t.Fatalf("Failed to append common flags, %v", err)
	}

	err = ValidateCommonFlags(fs)

	if err != nil {
		t.Fatalf("Failed to validate common flags, %v", err)
	}

}
