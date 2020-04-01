package app

import (
	_ "github.com/whosonfirst/go-reader-http"
	_ "github.com/whosonfirst/go-whosonfirst-spatial/database/extras"
)

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
)

func NewExtrasDatabaseWithFlagSet(ctx context.Context, fl *flag.FlagSet) (database.ExtrasDatabase, error) {

	enable_extras, _ := flags.BoolVar(fl, "enable-extras")
	extras_database_uri, _ := flags.StringVar(fl, "extras-database")

	if !enable_extras {
		return nil, nil
	}

	return database.NewExtrasDatabase(ctx, extras_database_uri)
}
