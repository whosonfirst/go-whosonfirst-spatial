package app

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-spatial/flags"
	"github.com/whosonfirst/go-spatial/index"
	wof_index "github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"runtime/debug"
	"time"
)

type PIPApplication struct {
	mode   string
	Index  index.Index
	Cache  cache.Cache
	Extras *database.SQLiteDatabase
	Walker *wof_index.Indexer
	Logger *log.WOFLogger
}

func NewPIPApplication(ctx context.Context, fl *flag.FlagSet) (*PIPApplication, error) {

	logger, err := NewApplicationLogger(ctx, fl)

	if err != nil {
		return nil, err
	}

	appcache, err := NewApplicationCache(ctx, fl)

	if err != nil {
		return nil, err
	}

	appindex, err := NewApplicationIndex(ctx, fl)

	if err != nil {
		return nil, err
	}

	appextras, err := NewApplicationExtras(ctx, fl)

	if err != nil {
		return nil, err
	}

	walker, err := NewApplicationWalker(ctx, fl, appindex, appextras)

	if err != nil {
		return nil, err
	}

	mode, _ := flags.StringVar(fl, "mode")

	p := PIPApplication{
		mode:   mode,
		Cache:  appcache,
		Index:  appindex,
		Extras: appextras,
		Walker: walker,
		Logger: logger,
	}

	return &p, nil
}

func (p *PIPApplication) Close(ctx context.Context) error {

	p.Cache.Close(ctx)
	p.Index.Close(ctx)

	if p.Extras != nil {
		p.Extras.Close()
	}

	return nil
}

func (p *PIPApplication) IndexPaths(paths []string) error {

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
