package alternate

import (
	"github.com/whosonfirst/go-whosonfirst-flags"
	// "github.com/whosonfirst/go-whosonfirst-sources"
)

type AlternateGeometryFlag struct {
	flags.AlternateGeometryFlag
	label string
}

func NewAlternateGeometryFlag(label string) (flags.AlternateGeometryFlag, error) {

	// check label against go-whosonfirst-sources here?

	f := AlternateGeometryFlag{
		label: label,
	}

	return &f, nil
}

func (f *AlternateGeometryFlag) MatchesAny(others ...flags.AlternateGeometryFlag) bool {

	for _, o := range others {

		if f.Label() == o.Label() {
			return true
		}

	}

	return false
}

func (f *AlternateGeometryFlag) MatchesAll(others ...flags.AlternateGeometryFlag) bool {

	matches := 0

	for _, o := range others {

		if f.Label() == o.Label() {
			matches += 1
		}

	}

	if matches == len(others) {
		return true
	}

	return false
}

func (f *AlternateGeometryFlag) Label() string {
	return f.label
}

func (f *AlternateGeometryFlag) String() string {
	return f.Label()
}
