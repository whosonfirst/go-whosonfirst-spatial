package app

import (
	_ "github.com/whosonfirst/go-reader-http"
)

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-spatial/extras"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
)

func NewExtrasReaderWithFlagSet(ctx context.Context, fl *flag.FlagSet) (extras.ExtrasReader, error) {

	enable_extras, _ := flags.BoolVar(fl, "enable-extras")
	extras_reader_uri, _ := flags.StringVar(fl, "extras-reader-uri")

	if !enable_extras {
		return nil, nil
	}

	return extras.NewExtrasReader(ctx, extras_reader_uri)
}
