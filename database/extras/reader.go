package extras

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-reader-cachereader"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/geojson"
	_ "log"
	"net/url"
)

func init() {
	ctx := context.Background()
	database.RegisterExtrasDatabase(ctx, "reader", NewReaderExtrasDatabase)
}

type ExtrasResponse struct {
	Index   int
	Feature geojson.GeoJSONFeature
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

func (db *ReaderExtrasDatabase) IndexFeature(context.Context, wof_geojson.Feature) error {
	return nil
}

func (db *ReaderExtrasDatabase) AppendExtrasWithFeatureCollection(ctx context.Context, fc *geojson.GeoJSONFeatureCollection, extras []string) (*geojson.GeoJSONFeatureCollection, error) {

	rsp_ch := make(chan ExtrasResponse)
	err_ch := make(chan error)
	done_ch := make(chan bool)

	remaining := len(fc.Features)

	for idx, f := range fc.Features {
		go db.appendExtrasWithChannels(ctx, idx, f, extras, rsp_ch, err_ch, done_ch)
	}

	for remaining > 0 {
		select {
		case <-ctx.Done():
			return nil, nil
		case <-done_ch:
			remaining -= 1
		case rsp := <-rsp_ch:
			fc.Features[rsp.Index] = rsp.Feature
		case err := <-err_ch:
			return nil, err
		default:
			// pass
		}
	}

	return fc, nil
}

func (db *ReaderExtrasDatabase) appendExtrasWithChannels(ctx context.Context, idx int, f geojson.GeoJSONFeature, extras []string, rsp_ch chan ExtrasResponse, err_ch chan error, done_ch chan bool) {

	defer func() {
		done_ch <- true
	}()

	select {
	case <-ctx.Done():
		return
	default:
		// pass
	}

	target, err := json.Marshal(f)

	if err != nil {
		err_ch <- err
		return
	}

	id_rsp := gjson.GetBytes(target, "properties.wof:id")

	if !id_rsp.Exists() {
		err_ch <- errors.New("Missing wof:id")
		return
	}

	id := id_rsp.Int()

	source, err := wof_reader.LoadBytesFromID(ctx, db.reader, id)

	if err != nil {
		err_ch <- err
		return
	}

	target, err = AppendExtrasWithBytes(ctx, source, target, extras)

	if err != nil {
		err_ch <- err
		return
	}

	var new_f geojson.GeoJSONFeature
	err = json.Unmarshal(target, &new_f)

	if err != nil {
		err_ch <- err
		return
	}

	rsp := ExtrasResponse{
		Index:   idx,
		Feature: new_f,
	}

	rsp_ch <- rsp
	return
}
