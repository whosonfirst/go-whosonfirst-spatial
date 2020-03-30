package main

// go run -mod vendor cmd/wof-pip-server/main.go -index 'rtree://' -mode repo:// /usr/local/data/sfomuseum-data-maps/

import (
	"context"
	"fmt"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/assets/templates"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial/http"
	"html/template"
	"log"
	gohttp "net/http"
	"os"
)

func main() {

	fs, err := flags.CommonFlags()

	if err != nil {
		log.Fatal(err)
	}

	err = flags.AppendWWWFlags(fs)

	flags.Parse(fs)

	ctx := context.Background()

	err = flags.ValidateCommonFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	err = flags.ValidateWWWFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	pip, err := app.NewSpatialApplicationWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create new PIP application, because %s", err))
	}

	err = pip.IndexPaths(fs.Args())

	if err != nil {
		pip.Logger.Fatal("Failed to index paths, because %s", err)
	}

	pip.Logger.Debug("setting up intersects handler")

	enable_geojson, _ := flags.BoolVar(fs, "enable-geojson")

	intersects_opts := http.NewDefaultIntersectsHandlerOptions()
	intersects_opts.EnableGeoJSON = enable_geojson

	intersects_handler, err := http.IntersectsHandler(pip, intersects_opts)

	if err != nil {
		pip.Logger.Fatal("failed to create intersects handler because %s", err)
	}

	ping_handler, err := http.PingHandler()

	if err != nil {
		pip.Logger.Fatal("failed to create ping handler because %s", err)
	}

	mux := gohttp.NewServeMux()

	mux.Handle("/ping", ping_handler)
	mux.Handle("/intersects", intersects_handler)

	enable_www, _ := flags.BoolVar(fs, "enable-www")
	enable_candidates, _ := flags.BoolVar(fs, "enable-candidates")

	if enable_candidates {

		pip.Logger.Debug("setting up candidates handler")

		candidateshandler, err := http.IntersectsCandidatesHandler(pip)

		if err != nil {
			pip.Logger.Fatal("failed to create Spatial handler because %s", err)
		}

		mux.Handle("/intersects/candidates", candidateshandler)
	}

	if enable_www {

		path_templates, _ := flags.StringVar(fs, "path-templates")

		t := template.New("spatial").Funcs(template.FuncMap{
			//
		})

		if path_templates != "" {

			t, err = t.ParseGlob(path_templates)

			if err != nil {
				log.Fatal(err)
			}

		} else {

			for _, name := range templates.AssetNames() {

				body, err := templates.Asset(name)

				if err != nil {
					log.Fatal(err)
				}

				t, err = t.Parse(string(body))

				if err != nil {
					log.Fatal(err)
				}
			}
		}

		intersects_opts.Templates = t

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

		www_handler, err := http.IntersectsWWWHandler(pip, intersects_opts)

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
