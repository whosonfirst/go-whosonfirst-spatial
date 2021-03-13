package filter

import (
	"flag"
	"github.com/sfomuseum/go-flags/lookup"
)

func NewSPRFilterFromFlagSet(fs *flag.FlagSet) (Filter, error) {

	inputs, err := NewSPRInputsFromFlagSet(fs)

	if err != nil {
		return nil, err
	}

	return NewSPRFilterFromInputs(inputs)
}

func NewSPRInputsFromFlagSet(fs *flag.FlagSet) (*SPRInputs, error) {

	inputs, err := NewSPRInputs()

	if err != nil {
		return nil, err
	}

	placetypes, err := lookup.MultiStringVar(fs, "placetype")

	if err != nil {
		return nil, err
	}

	inputs.Placetypes = placetypes

	inception_date, err := lookup.StringVar(fs, "inception-date")

	if err != nil {
		return nil, err
	}

	inputs.InceptionDate = inception_date

	cessation_date, err := lookup.StringVar(fs, "cessation-date")

	if err != nil {
		return nil, err
	}

	inputs.CessationDate = cessation_date

	geometries, err := lookup.MultiStringVar(fs, "geometries")

	if err != nil {
		return nil, err
	}

	inputs.Geometries = geometries

	alt_geoms, err := lookup.MultiStringVar(fs, "alternate-geometry")

	if err != nil {
		return nil, err
	}

	inputs.AlternateGeometries = alt_geoms

	is_current, err := lookup.MultiInt64Var(fs, "is-current")

	if err != nil {
		return nil, err
	}

	inputs.IsCurrent = is_current

	is_ceased, err := lookup.MultiInt64Var(fs, "is-ceased")

	if err != nil {
		return nil, err
	}

	inputs.IsCeased = is_ceased

	is_deprecated, err := lookup.MultiInt64Var(fs, "is-deprecated")

	if err != nil {
		return nil, err
	}

	inputs.IsDeprecated = is_deprecated

	is_superseded, err := lookup.MultiInt64Var(fs, "is-superseded")

	if err != nil {
		return nil, err
	}

	inputs.IsSuperseded = is_superseded

	is_superseding, err := lookup.MultiInt64Var(fs, "is-superseding")

	if err != nil {
		return nil, err
	}

	inputs.IsSuperseding = is_superseding

	return inputs, nil
}
