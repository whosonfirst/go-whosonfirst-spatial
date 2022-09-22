package pointinpolygon

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-spatial/api"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"log"
)

func Run(ctx context.Context, logger *log.Logger) error {
	fs, err := DefaultFlagSet()

	if err != nil {
		return fmt.Errorf("Failed to derive flags, %w", err)
	}

	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	err := flags.ValidateCommonFlags(fs)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = flags.ValidateQueryFlags(fs)

	if err != nil {
		return fmt.Errorf("%w", err)
	}

	database_uri, _ := lookup.StringVar(fs, flags.SPATIAL_DATABASE_URI)

	db, err := database.NewSpatialDatabase(ctx, database_uri)

	if err != nil {
		return fmt.Errorf("Failed to create database for '%s', %v", database_uri, err)
	}

	query := func(ctx context.Context, req *api.PointInPolygonRequest) (interface{}, error) {

		c, err := geo.NewCoordinate(req.Longitude, req.Latitude)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new coordinate, %v", err)
		}

		f, err := api.NewSPRFilterFromPointInPolygonRequest(req)

		if err != nil {
			return nil, err
		}

		r, err := db.PointInPolygon(ctx, c, f)

		if err != nil {
			return nil, fmt.Errorf("Failed to query database with coord %v, %v", c, err)
		}

		return r, nil
	}

	req, err := api.NewPointInPolygonRequestFromFlagSet(fs)

	if err != nil {
		return fmt.Errorf("Failed to create SPR filter, %v", err)
	}

	rsp, err := query(ctx, req)

	if err != nil {
		return fmt.Errorf("Failed to query, %v", err)
	}

	enc, err := json.Marshal(rsp)

	if err != nil {
		return fmt.Errorf("Failed to marshal results, %v", err)
	}

	fmt.Println(string(enc))
	return nil
}
