package app

import (
	"errors"
	"flag"
	"github.com/whosonfirst/go-spatial/flags"
	"github.com/whosonfirst/go-spatial/index"
)

func NewApplicationIndex(fl *flag.FlagSet) (index.Index, error) {

	index_uri, err := flags.StringVar(fl, "index")

	if err != nil {
		return nil, err
	}

	return index.NewIndex(index_uri)
}
