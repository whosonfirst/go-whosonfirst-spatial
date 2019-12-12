package app

import (
	"errors"
	"flag"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-spatial/flags"
)

func NewApplicationCache(fl *flag.FlagSet) (cache.Cache, error) {

	cache_uri, err := flags.StringVar(fl, "cache")

	if err != nil {
		return nil, err
	}

	return cache.NewCache(cache_uri)
}
