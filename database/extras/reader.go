package extras

import (
	"context"
	"errors"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spr"		
	"github.com/whosonfirst/go-reader-cachereader"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-cache"	
	"net/url"
)

func init() {
	ctx := context.Background()
	database.RegisterExtrasDatabase(ctx, "reader", NewReaderExtrasDatabase)
}

type ReaderExtrasDatabase struct {
	database.ExtrasDatabase
	reader reader.Reader
}

func NewReaderExtrasDatabase(ctx context.Context, uri string) (database.ExtrasDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	reader_uri := q.Get("reader")

	if reader_uri == "" {
		return nil, errors.New("Missing reader parameter")
	}
	
	cache_uri := q.Get("cache")

	if cache_uri == "" {
		cache_uri = "null://"
	}

	r, err := reader.NewReader(ctx, reader_uri)

	if err != nil {
		return nil, err
	}

	c, err := cache.NewCache(ctx, cache_uri)

	if err != nil {
		return nil, err
	}

	cr, err := cachereader.NewCacheReader(r, c)

	if err != nil {
		return nil, err
	}
	
	db := &ReaderExtrasDatabase{
		reader: cr,
	}

	return db, nil
}

func (db *ReaderExtrasDatabase) Close(ctx context.Context) error {
	return nil
}

func (db *ReaderExtrasDatabase) IndexFeature(context.Context, geojson.Feature) error {
	return nil
}

func (db *ReaderExtrasDatabase) AppendExtrasWithSPRResults(context.Context, spr.StandardPlacesResults, ...string) error {
	return nil
}
