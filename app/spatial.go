package app

import (
	_ "github.com/whosonfirst/go-whosonfirst-spatial/database/spatial"
)

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
)

func NewSpatialDatabase(ctx context.Context, fl *flag.FlagSet) (database.SpatialDatabase, error) {

	index_uri, err := flags.StringVar(fl, "index")

	if err != nil {
		return nil, err
	}

	return database.NewSpatialDatabase(ctx, index_uri)
}
