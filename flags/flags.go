package flags

import (
	"errors"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/whosonfirst/go-whosonfirst-index"
	"log"
	"os"
	"sort"
	"strings"
)

func Parse(fs *flag.FlagSet) {

	args := os.Args[1:]

	if len(args) > 0 && args[0] == "-h" {
		fs.Usage()
		os.Exit(0)
	}

	if len(args) > 0 && args[0] == "-setenv" {
		SetFromEnv(fs)
	}

	fs.Parse(args)
}

func SetFromEnv(fs *flag.FlagSet) {

	fs.VisitAll(func(fl *flag.Flag) {

		name := fl.Name
		env := name

		env = strings.ToUpper(env)
		env = strings.Replace(env, "-", "_", -1)
		env = fmt.Sprintf("WOF_%s", env)

		val, ok := os.LookupEnv(env)

		if ok {
			log.Printf("set -%s flag (%s) from %s environment variable\n", name, val, env)
			fs.Set(name, val)
		}

	})
}

func ValidateCommonFlags(fs *flag.FlagSet) error {

	mode, err := StringVar(fs, "mode")

	if err != nil {
		return err
	}

	pip_index, err := StringVar(fs, "index")

	if err != nil {
		return err
	}

	pip_cache, err := StringVar(fs, "cache")

	if err != nil {
		return err
	}

	spatialite_dsn, err := StringVar(fs, "spatialite-dsn")

	if err != nil {
		return err
	}

	if mode == "spatialite" {

		if pip_index != "spatialite" {
			return errors.New("-mode is spatialite but -index is not")
		}

		if pip_cache != "sqlite" && pip_cache != "spatialite" {
			return errors.New("-mode is spatialite but -cache is neither 'sqlite' or 'spatialite'")
		}

		if spatialite_dsn == "" || spatialite_dsn == ":memory:" {
			return errors.New("-spatialite-dsn needs to be an actual file on disk")
		}

		enable_extras, err := BoolVar(fs, "enable-extras")

		if err != nil {
			return err
		}

		if enable_extras {

			extras_dsn, err := StringVar(fs, "extras-dsn")

			if err != nil {
				return err
			}

			if extras_dsn == ":tmpfile:" {
				log.Println("-mode is spatialite so assigning the value of -spatialite-dsn to -extras-dsn")
				fs.Set("extras-dsn", spatialite_dsn)
			} else if extras_dsn != spatialite_dsn {
				return errors.New("-mode is spatialite so -extras-dsn needs to be the same as -spatialite-dsn")
			} else {
				// pass
			}
		}
	}

	return nil
}

func ValidateWWWFlags(fs *flag.FlagSet) error {

	strict, err := BoolVar(fs, "strict")

	if err != nil {
		return err
	}

	enable_www, err := BoolVar(fs, "enable-www")

	if err != nil {
		return err
	}

	if enable_www {

		log.Println("-enable-www flag is true causing the following flags to also be true: -enable-geojson -enable-candidates")

		fs.Set("enable-geojson", "true")
		fs.Set("enable-candidates", "true")

		key, err := StringVar(fs, "www-api-key")

		if err != nil {
			return err
		}

		if key == "xxxxxx" {

			warning := "-enable-www flag is set but -www-api-key is empty"

			if strict {
				return errors.New(warning)
			}

			log.Printf("[WARNING] %s\n", warning)
		}
	}

	return nil
}

func Lookup(fl *flag.FlagSet, k string) (interface{}, error) {

	v := fl.Lookup(k)

	if v == nil {
		msg := fmt.Sprintf("Unknown flag '%s'", k)
		return nil, errors.New(msg)
	}

	// Go is weird...
	return v.Value.(flag.Getter).Get(), nil
}

func StringVar(fl *flag.FlagSet, k string) (string, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return "", err
	}

	return i.(string), nil
}

func IntVar(fl *flag.FlagSet, k string) (int, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return 0, err
	}

	return i.(int), nil
}

func BoolVar(fl *flag.FlagSet, k string) (bool, error) {

	i, err := Lookup(fl, k)

	if err != nil {
		return false, err
	}

	return i.(bool), nil
}

func NewFlagSet(name string) *flag.FlagSet {

	fs := flag.NewFlagSet(name, flag.ExitOnError)

	fs.Usage = func() {
		fs.PrintDefaults()
	}

	return fs
}

func CommonFlags() (*flag.FlagSet, error) {

	fs := NewFlagSet("common")

	fs.String("index", "rtree", "Valid options are: rtree, spatialite.")
	fs.String("cache", "gocache", "Valid options are: gocache, fs, spatialite, sqlite. Note that the spatalite option is just a convenience to mirror the '-index spatialite' option.")

	modes := index.Modes()
	modes = append(modes, "spatialite")

	sort.Strings(modes)

	valid_modes := strings.Join(modes, ", ")
	desc_modes := fmt.Sprintf("Valid modes are: %s.", valid_modes)

	fs.String("mode", "files", desc_modes)

	fs.String("spatialite-dsn", "", "A valid SQLite DSN for the '-cache spatialite/sqlite' or '-index spatialite' option. As of this writing for the '-index' and '-cache' options share the same '-spatailite' DSN.")
	fs.String("fs-path", "", "The root directory to look for features if '-cache fs'.")

	fs.Bool("is-wof", true, "Input data is WOF-flavoured GeoJSON. (Pass a value of '0' or 'false' if you need to index non-WOF documents.")

	// this is invoked/used in app/indexer.go but for the life of me I can't
	// figure out how to make the code in flags/exclude.go implement the
	// correct inferface wah wah so that flag.Lookup("exclude").Value returns
	// something we can loop over... so instead we just strings.Split() on
	// flag.Lookup("exclude").String() which is dumb but works...
	// (20180301/thisisaaronland)

	var exclude Exclude
	fs.Var(&exclude, "exclude", "Exclude (WOF) records based on their existential flags. Valid options are: ceased, deprecated, not-current, superseded.")

	fs.Bool("setenv", false, "Set flags from environment variables.")
	fs.Bool("verbose", false, "Be chatty.")
	fs.Bool("strict", false, "Be strict about flags and fail if any are missing or deprecated flags are used.")

	return fs, nil
}

func AppendWWWFlags(fs *flag.FlagSet) error {

	fs.String("host", "localhost", "The hostname to listen for requests on.")
	fs.Int("port", 8080, "The port number to listen for requests on.")

	fs.Bool("enable-extras", false, "Enable support for 'extras' parameters in queries.")
	fs.String("extras-dsn", ":tmpfile:", "A valid SQLite DSN for your 'extras' database - if ':tmpfile:' then a temporary database will be created during indexing and deleted when the program exits.")

	fs.Bool("enable-geojson", false, "Allow users to request GeoJSON FeatureCollection formatted responses.")
	fs.Bool("enable-candidates", false, "Enable the /candidates endpoint to return candidate bounding boxes (as GeoJSON) for requests.")
	fs.Bool("enable-polylines", false, "Enable the /polylines endpoint to return hierarchies intersecting a path.")
	fs.Bool("enable-www", false, "Enable the interactive /debug endpoint to query points and display results.")

	fs.Int("polylines-max-coords", 100, "The maximum number of points a (/polylines) path may contain before it is automatically paginated.")
	fs.String("www-path", "/debug", "The URL path for the interactive debug endpoint.")
	fs.String("www-api-key", "xxxxxx", "A valid Nextzen Map Tiles API key (https://developers.nextzen.org).")

	fs.String("static-prefix", "", "Prepend this prefix to URLs for static assets.")

	fs.String("nextzen-apikey", "", "A valid Nextzen API key")
	fs.String("nextzen-style-url", "/tangram/refill-style.zip", "...")
	fs.String("nextzen-tile-url", tangramjs.NEXTZEN_MVT_ENDPOINT, "...")

	fs.String("templates", "", "An optional string for local templates. This is anything that can be read by the 'templates.ParseGlob' method.")

	return nil
}
