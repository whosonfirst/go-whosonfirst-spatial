package app

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-spatial/flags"
)

func NewApplicationCache(ctx context.Context, fl *flag.FlagSet) (cache.Cache, error) {

	cache_uri, err := flags.StringVar(fl, "cache")

	if err != nil {
		return nil, err
	}

	return cache.NewCache(ctx, cache_uri)
}
