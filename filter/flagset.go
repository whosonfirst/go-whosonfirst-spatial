package filter

import (
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
)

func NewSPRFilterFromFlagSet(fs *flag.FlagSet) (spatial.Filter, error) {

	inputs, err := NewSPRInputsFromFlagSet(fs)

	if err != nil {
		return nil, fmt.Errorf("Failed to create SPR inputs from flagset, %w", err)
	}

	return NewSPRFilterFromInputs(inputs)
}

func NewSPRInputsFromFlagSet(fs *flag.FlagSet) (*SPRInputs, error) {

	inputs, err := NewSPRInputs()

	if err != nil {
		return nil, fmt.Errorf("Failed to create SPR inputs, %w", err)
	}

	placetypes, err := lookup.MultiStringVar(fs, flags.PlacetypeFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup %s flag, %w", flags.PlacetypeFlag, err)
	}

	inputs.Placetypes = placetypes

	inception_date, err := lookup.StringVar(fs, flags.InceptionDateFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup %s flag, %w", flags.InceptionDateFlag, err)
	}

	inputs.InceptionDate = inception_date

	cessation_date, err := lookup.StringVar(fs, flags.CessationDateFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup %s flag, %w", flags.CessationDateFlag, err)
	}

	inputs.CessationDate = cessation_date

	geometries, err := lookup.StringVar(fs, flags.GeometriesFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup %s flag, %w", flags.GeometriesFlag, err)
	}

	inputs.Geometries = []string{geometries}

	alt_geoms, err := lookup.MultiStringVar(fs, flags.AlternateGeometriesFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup %s flag, %w", flags.AlternateGeometriesFlag, err)
	}

	inputs.AlternateGeometries = alt_geoms

	is_current, err := lookup.MultiInt64Var(fs, flags.IsCurrentFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup %s flag, %w", flags.IsCurrentFlag, err)
	}

	inputs.IsCurrent = is_current

	is_ceased, err := lookup.MultiInt64Var(fs, flags.IsCeasedFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup %s flag, %w", flags.IsCeasedFlag, err)
	}

	inputs.IsCeased = is_ceased

	is_deprecated, err := lookup.MultiInt64Var(fs, flags.IsDeprecatedFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup %s flag, %w", flags.IsDeprecatedFlag, err)
	}

	inputs.IsDeprecated = is_deprecated

	is_superseded, err := lookup.MultiInt64Var(fs, flags.IsSupersededFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup %s flag, %w", flags.IsSupersededFlag, err)
	}

	inputs.IsSuperseded = is_superseded

	is_superseding, err := lookup.MultiInt64Var(fs, flags.IsSupersedingFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup %s flag, %w", flags.IsSupersedingFlag, err)
	}

	inputs.IsSuperseding = is_superseding

	return inputs, nil
}
