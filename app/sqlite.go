package app

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
)

func NewSpatialiteDB(ctx context.Context, fl *flag.FlagSet) (*database.SQLiteDatabase, error) {

	dsn, err := flags.StringVar(fl, "spatialite-dsn")

	if err != nil {
		return nil, err
	}

	db, err := database.NewDBWithDriver("spatialite", dsn)

	if err != nil {
		return nil, err
	}

	err = db.LiveHardDieFast()

	if err != nil {
		return nil, err
	}

	return db, nil
}
