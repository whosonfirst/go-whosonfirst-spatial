package app

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-cache"
	wof_index "github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"runtime/debug"
	"time"
)

type SpatialApplication struct {
	mode            string
	SpatialDatabase database.SpatialDatabase
	Cache           cache.Cache
	ExtrasDatabase  database.ExtrasDatabase
	Walker          *wof_index.Indexer
	Logger          *log.WOFLogger
}

func NewSpatialApplication(ctx context.Context, fl *flag.FlagSet) (*SpatialApplication, error) {

	logger, err := NewApplicationLogger(ctx, fl)

	if err != nil {
		return nil, err
	}

	appcache, err := NewCache(ctx, fl)

	if err != nil {
		return nil, err
	}

	spatial_db, err := NewSpatialDatabase(ctx, fl)

	if err != nil {
		return nil, err
	}

	extras_db, err := NewExtrasDatabase(ctx, fl)

	if err != nil {
		return nil, err
	}

	walker, err := NewWalker(ctx, fl, spatial_db, extras_db)

	if err != nil {
		return nil, err
	}

	mode, _ := flags.StringVar(fl, "mode")

	sp := SpatialApplication{
		mode:            mode,
		Cache:           appcache,
		SpatialDatabase: spatial_db,
		ExtrasDatabase:  extras_db,
		Walker:          walker,
		Logger:          logger,
	}

	return &sp, nil
}

func (p *SpatialApplication) Close(ctx context.Context) error {

	p.Cache.Close(ctx)
	p.SpatialDatabase.Close(ctx)

	if p.ExtrasDatabase != nil {
		p.ExtrasDatabase.Close(ctx)
	}

	return nil
}

func (p *SpatialApplication) IndexPaths(paths []string) error {

	if p.mode != "spatialite" {

		go func() {

			// TO DO: put this somewhere so that it can be triggered by signal(s)
			// to reindex everything in bulk or incrementally

			t1 := time.Now()

			err := p.Walker.IndexPaths(paths)

			if err != nil {
				p.Logger.Fatal("failed to index paths because %s", err)
			}

			t2 := time.Since(t1)

			p.Logger.Status("finished indexing in %v", t2)
			debug.FreeOSMemory()
		}()

		// set up some basic monitoring and feedback stuff

		go func() {

			c := time.Tick(1 * time.Second)

			for _ = range c {

				if !p.Walker.IsIndexing() {
					continue
				}

				p.Logger.Status("indexing %d records indexed", p.Walker.Indexed)
			}
		}()
	}

	return nil
}
