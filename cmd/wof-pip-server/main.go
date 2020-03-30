package main

import (
	"context"
	"fmt"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial/http"
	"log"
	gohttp "net/http"
	"os"
	"runtime"
	"time"
)

func main() {

	fs, err := flags.CommonFlags()

	if err != nil {
		log.Fatal(err)
	}

	err = flags.AppendWWWFlags(fs)

	flags.Parse(fs)

	err = flags.ValidateCommonFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	err = flags.ValidateWWWFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pip, err := app.NewPIPApplication(ctx, fs)

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create new PIP application, because %s", err))
	}

	pip_index, _ := flags.StringVar(fs, "index")
	pip_cache, _ := flags.StringVar(fs, "cache")
	mode, _ := flags.StringVar(fs, "mode")

	pip.Logger.Info("index is %s cache is %s mode is %s", pip_index, pip_cache, mode)

	err = pip.IndexPaths(fs.Args())

	if err != nil {
		pip.Logger.Fatal("Failed to index paths, because %s", err)
	}

	go func() {

		tick := time.Tick(1 * time.Minute)

		for _ = range tick {
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			pip.Logger.Status("memstats system: %8d inuse: %8d released: %8d objects: %6d", ms.HeapSys, ms.HeapInuse, ms.HeapReleased, ms.HeapObjects)
		}
	}()

	// set up the HTTP endpoint

	pip.Logger.Debug("setting up intersects handler")

	enable_geojson, _ := flags.BoolVar(fs, "enable-geojson")

	intersects_opts := http.NewDefaultIntersectsHandlerOptions()
	intersects_opts.EnableGeoJSON = enable_geojson

	intersects_handler, err := http.IntersectsHandler(pip.Index, pip.Walker, pip.Extras, intersects_opts)

	if err != nil {
		pip.Logger.Fatal("failed to create PIP handler because %s", err)
	}

	ping_handler, err := http.PingHandler()

	if err != nil {
		pip.Logger.Fatal("failed to create Ping handler because %s", err)
	}

	mux := gohttp.NewServeMux()

	mux.Handle("/ping", ping_handler)
	mux.Handle("/intersects", intersects_handler)

	enable_www, _ := flags.BoolVar(fs, "enable-www")
	enable_candidates, _ := flags.BoolVar(fs, "enable-candidates")

	if enable_candidates {

		pip.Logger.Debug("setting up candidates handler")

		candidateshandler, err := http.IntersectsCandidatesHandler(pip.Index, pip.Walker)

		if err != nil {
			pip.Logger.Fatal("failed to create Spatial handler because %s", err)
		}

		mux.Handle("/intersects/candidates", candidateshandler)
	}

	if enable_www {

		nextzen_apikey, _ := flags.StringVar(fs, "nextzen-apikey")
		nextzen_style_url, _ := flags.StringVar(fs, "nextzen-style-url")
		nextzen_tile_url, _ := flags.StringVar(fs, "nextzen-tile-url")

		static_prefix, _ := flags.StringVar(fs, "static-prefix")

		bootstrap_opts := bootstrap.DefaultBootstrapOptions()

		tangramjs_opts := tangramjs.DefaultTangramJSOptions()
		tangramjs_opts.Nextzen.APIKey = nextzen_apikey
		tangramjs_opts.Nextzen.StyleURL = nextzen_style_url
		tangramjs_opts.Nextzen.TileURL = nextzen_tile_url

		err = bootstrap.AppendAssetHandlersWithPrefix(mux, static_prefix)

		www_path, _ := flags.StringVar(fs, "www-path")
		www_handler, err := http.BundledWWWHandler()

		if err != nil {
			pip.Logger.Fatal("failed to create (bundled) www handler because %s", err)
		}

		www_handler = bootstrap.AppendResourcesHandlerWithPrefix(www_handler, bootstrap_opts, static_prefix)
		www_handler = tangramjs.AppendResourcesHandlerWithPrefix(www_handler, tangramjs_opts, static_prefix)

		mux.Handle(www_path, www_handler)
	}

	host, _ := flags.StringVar(fs, "host")
	port, _ := flags.IntVar(fs, "port")

	endpoint := fmt.Sprintf("%s:%d", host, port)
	pip.Logger.Status("listening for requests on %s", endpoint)

	err = gohttp.ListenAndServe(endpoint, mux)

	if err != nil {
		pip.Logger.Fatal("failed to start server because %s", err)
	}

	os.Exit(0)
}
