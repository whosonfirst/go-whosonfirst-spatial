package main

// go run -mod vendor cmd/spatial-server/main.go -index 'rtree://' -mode repo:// /usr/local/data/sfomuseum-data-maps/

import (
	"context"
	"fmt"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/assets/templates"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial/http"
	"github.com/whosonfirst/go-whosonfirst-spatial/server"
	"html/template"
	"log"
	gohttp "net/http"
	gourl "net/url"
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

	spatial_app, err := app.NewSpatialApplicationWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create new spatial application, because %s", err))
	}

	logger := spatial_app.Logger

	paths := fs.Args()

	err = spatial_app.IndexPaths(ctx, paths...)

	if err != nil {
		logger.Fatal("Failed to index paths, because %s", err)
	}

	logger.Debug("setting up intersects handler")

	enable_geojson, _ := flags.BoolVar(fs, "enable-geojson")

	intersects_opts := &http.IntersectsHandlerOptions{
		EnableGeoJSON: enable_geojson,
	}

	intersects_handler, err := http.IntersectsHandler(spatial_app, intersects_opts)

	if err != nil {
		logger.Fatal("failed to create intersects handler because %s", err)
	}

	ping_handler, err := http.PingHandler()

	if err != nil {
		logger.Fatal("failed to create ping handler because %s", err)
	}

	mux := gohttp.NewServeMux()

	mux.Handle("/ping", ping_handler)
	mux.Handle("/intersects", intersects_handler)

	enable_www, _ := flags.BoolVar(fs, "enable-www")
	enable_candidates, _ := flags.BoolVar(fs, "enable-candidates")

	if enable_candidates {

		logger.Debug("setting up candidates handler")

		candidateshandler, err := http.IntersectsCandidatesHandler(spatial_app)

		if err != nil {
			logger.Fatal("failed to create Spatial handler because %s", err)
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
				logger.Fatal("Unable to parse templates, %v", err)
			}

		} else {

			for _, name := range templates.AssetNames() {

				body, err := templates.Asset(name)

				if err != nil {
					logger.Fatal("Unable to load template '%s', %v", name, err)
				}

				t, err = t.Parse(string(body))

				if err != nil {
					logger.Fatal("Unable to parse template '%s', %v", name, err)
				}
			}
		}

		nextzen_apikey, _ := flags.StringVar(fs, "nextzen-apikey")
		nextzen_style_url, _ := flags.StringVar(fs, "nextzen-style-url")
		nextzen_tile_url, _ := flags.StringVar(fs, "nextzen-tile-url")

		bootstrap_opts := bootstrap.DefaultBootstrapOptions()

		tangramjs_opts := tangramjs.DefaultTangramJSOptions()
		tangramjs_opts.Nextzen.APIKey = nextzen_apikey
		tangramjs_opts.Nextzen.StyleURL = nextzen_style_url
		tangramjs_opts.Nextzen.TileURL = nextzen_tile_url

		err = tangramjs.AppendAssetHandlers(mux)

		if err != nil {
			logger.Fatal("Failed to append tangram.js assets, %v", err)
		}

		err = bootstrap.AppendAssetHandlers(mux)

		if err != nil {
			logger.Fatal("Failed to append bootstrap assets, %v", err)
		}

		intersects_www_opts := &http.IntersectsWWWHandlerOptions{
			Templates: t,
		}

		intersects_www_handler, err := http.IntersectsWWWHandler(spatial_app, intersects_www_opts)

		if err != nil {
			logger.Fatal("failed to create (bundled) www handler because %s", err)
		}

		intersects_www_handler = bootstrap.AppendResourcesHandler(intersects_www_handler, bootstrap_opts)
		intersects_www_handler = tangramjs.AppendResourcesHandler(intersects_www_handler, tangramjs_opts)

		www_path, _ := flags.StringVar(fs, "www-path")
		mux.Handle(www_path, intersects_www_handler)
	}

	host, _ := flags.StringVar(fs, "host")
	port, _ := flags.IntVar(fs, "port")
	proto := "http" // FIX ME

	address := fmt.Sprintf("spatial://%s:%d", host, port)

	u, err := gourl.Parse(address)

	if err != nil {
		logger.Fatal("Failed to parse address '%s', %v", address, err)
	}

	s, err := server.NewStaticServer(proto, u)

	if err != nil {
		logger.Fatal("Failed to create new server for '%s' (%s), %v", u, proto, err)
	}

	logger.Info("Listening on %s", s.Address())

	err = s.ListenAndServe(mux)

	if err != nil {
		logger.Fatal("Failed to start server, %v", err)
	}
}
