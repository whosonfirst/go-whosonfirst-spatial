package main

// go run -mod vendor cmd/spatial-server/main.go -enable-www -mode repo:// /usr/local/data/sfomuseum-data-maps/

import (
	"context"
	"fmt"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/assets/templates"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial/http/api"
	"github.com/whosonfirst/go-whosonfirst-spatial/http/health"
	"github.com/whosonfirst/go-whosonfirst-spatial/http/www"
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

	enable_geojson, _ := flags.BoolVar(fs, "enable-geojson")
	enable_www, _ := flags.BoolVar(fs, "enable-www")
	enable_candidates, _ := flags.BoolVar(fs, "enable-candidates")

	path_templates, _ := flags.StringVar(fs, "path-templates")
	nextzen_apikey, _ := flags.StringVar(fs, "nextzen-apikey")
	nextzen_style_url, _ := flags.StringVar(fs, "nextzen-style-url")
	nextzen_tile_url, _ := flags.StringVar(fs, "nextzen-tile-url")

	host, _ := flags.StringVar(fs, "host")
	port, _ := flags.IntVar(fs, "port")
	proto := "http" // FIX ME

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

	mux := gohttp.NewServeMux()

	ping_handler, err := health.PingHandler()

	if err != nil {
		logger.Fatal("failed to create ping handler because %s", err)
	}

	mux.Handle("/health/ping", ping_handler)

	api_pip_opts := &api.PointInPolygonHandlerOptions{
		EnableGeoJSON: enable_geojson,
	}

	api_pip_handler, err := api.PointInPolygonHandler(spatial_app, api_pip_opts)

	if err != nil {
		logger.Fatal("failed to create point-in-polygon handler because %s", err)
	}

	mux.Handle("/api/point-in-polygon", api_pip_handler)

	if enable_candidates {

		logger.Debug("setting up candidates handler")

		candidateshandler, err := api.PointInPolygonCandidatesHandler(spatial_app)

		if err != nil {
			logger.Fatal("failed to create Spatial handler because %s", err)
		}

		mux.Handle("/api/point-in-polygon/candidates", candidateshandler)
	}

	if enable_www {

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

		err = www.AppendStaticAssetHandlers(mux)

		if err != nil {
			logger.Fatal("Failed to append static assets, %v", err)
		}

		www_pip_opts := &www.PointInPolygonHandlerOptions{
			Templates: t,
		}

		www_pip_handler, err := www.PointInPolygonHandler(spatial_app, www_pip_opts)

		if err != nil {
			logger.Fatal("failed to create (bundled) www handler because %s", err)
		}

		www_pip_handler = bootstrap.AppendResourcesHandler(www_pip_handler, bootstrap_opts)
		www_pip_handler = tangramjs.AppendResourcesHandler(www_pip_handler, tangramjs_opts)

		mux.Handle("/point-in-polygon", www_pip_handler)
	}

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