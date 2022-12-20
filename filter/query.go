package filter

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"net/url"
	"strconv"
)

// NewSPRFilterFromQuery returns a new `spatial.Filter` instance derived from values in 'query'.
func NewSPRFilterFromQuery(query url.Values) (spatial.Filter, error) {

	inputs, err := NewSPRInputs()

	if err != nil {
		return nil, fmt.Errorf("Failed to create SPR inputs, %w", err)
	}

	inputs.Placetypes = query[flags.PlacetypeFlag]
	inputs.Geometries = query[flags.GeometriesFlag]
	inputs.AlternateGeometries = query[flags.AlternateGeometriesFlag]

	inputs.InceptionDate = query.Get(flags.InceptionDateFlag)
	inputs.CessationDate = query.Get(flags.CessationDateFlag)

	is_current, err := atoi(query[flags.IsCurrentFlag])

	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s flag, %w", flags.IsCurrentFlag, err)
	}

	is_deprecated, err := atoi(query[flags.IsDeprecatedFlag])

	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s flag, %w", flags.IsDeprecatedFlag, err)
	}

	is_ceased, err := atoi(query[flags.IsCeasedFlag])

	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s flag, %w", flags.IsCeasedFlag, err)
	}

	is_superseded, err := atoi(query[flags.IsSupersededFlag])

	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s flag, %w", flags.IsSupersededFlag, err)
	}

	is_superseding, err := atoi(query[flags.IsSupersedingFlag])

	if err != nil {
		return nil, fmt.Errorf("Failed to parse %s flag, %w", flags.IsSupersedingFlag, err)
	}

	inputs.IsCurrent = is_current
	inputs.IsDeprecated = is_deprecated
	inputs.IsCeased = is_ceased
	inputs.IsSuperseded = is_superseded
	inputs.IsSuperseding = is_superseding

	return NewSPRFilterFromInputs(inputs)
}

func atoi(strings []string) ([]int64, error) {

	numbers := make([]int64, len(strings))

	for idx, str := range strings {

		i, err := strconv.ParseInt(str, 10, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse '%s', %w", str, err)
		}

		numbers[idx] = i
	}

	return numbers, nil
}
