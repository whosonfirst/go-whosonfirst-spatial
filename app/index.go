package app

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-spatial/flags"
	"github.com/whosonfirst/go-spatial/index"
)

func NewApplicationIndex(ctx context.Context, fl *flag.FlagSet) (index.Index, error) {

	index_uri, err := flags.StringVar(fl, "index")

	if err != nil {
		return nil, err
	}

	return index.NewIndex(ctx, index_uri)
}
