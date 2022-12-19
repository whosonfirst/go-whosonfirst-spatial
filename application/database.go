package application

import (
	"context"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
)

func NewSpatialDatabaseWithFlagSet(ctx context.Context, fl *flag.FlagSet) (spatial.SpatialDatabase, error) {

	spatial_uri, err := lookup.StringVar(fl, flags.SpatialDatabaseURIFlag)

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup flag for %s, %w", flags.SpatialDatabaseURIFlag, err)
	}

	return spatial.NewSpatialDatabase(ctx, spatial_uri)
}
