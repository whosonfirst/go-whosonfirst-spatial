package www

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"html/template"
	_ "log"
	gohttp "net/http"
)

type PointInPolygonHandlerOptions struct {
	Templates *template.Template
}

func PointInPolygonHandler(spatial_app *app.SpatialApplication, opts *PointInPolygonHandlerOptions) (gohttp.Handler, error) {

	t := opts.Templates.Lookup("pointinpolygon")

	if t == nil {
		return nil, errors.New("Missing pointinpolygon template")
	}

	walker := spatial_app.Walker

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		if walker.IsIndexing() {
			gohttp.Error(rsp, "indexing records", gohttp.StatusServiceUnavailable)
			return
		}

		// important if we're trying to use this in a Lambda/API Gateway context

		rsp.Header().Set("Content-Type", "text/html; charset=utf-8")

		err := t.Execute(rsp, nil)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
