package pointinpolygon

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
)

func DefaultFlagSet() (*flag.FlagSet, error) {

	fs, err := flags.CommonFlags()

	if err != nil {
		return nil, fmt.Errorf("Failed to assign common flags, %w", err)
	}

	err = flags.AppendQueryFlags(fs)

	if err != nil {
		return nil, fmt.Errorf("Failed to append query flags, %w", err)
	}

	return fs, nil
}
