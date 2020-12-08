package flags

import (
	"errors"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
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

func NewFlagSet(name string) *flag.FlagSet {

	fs := flag.NewFlagSet(name, flag.ExitOnError)

	fs.Usage = func() {
		fs.PrintDefaults()
	}

	return fs
}

func CommonFlags() (*flag.FlagSet, error) {

	fs := NewFlagSet("common")

	// spatial databases

	availables_databases := database.Schemes()
	desc_databases := fmt.Sprintf("Valid options are: %s", available_databases)

	fs.String("spatial-database-uri", "rtree://", desc_databases)

	// property readers

	fs.Bool("enable-properties", false, "Enable support for 'properties' parameters in queries.")

	availables_property_readers := properties.Schemes()
	desc_property_readers := fmt.Sprintf("Valid options are: %s", available_property_readers)

	fs.String("properties-reader-uri", "", desc_property_readers)

	// indexing modes

	modes := index.Modes()
	sort.Strings(modes)

	valid_modes := strings.Join(modes, ", ")
	desc_modes := fmt.Sprintf("Valid modes are: %s.", valid_modes)

	fs.String("mode", "repo://", desc_modes)

	//

	fs.Bool("is-wof", true, "Input data is WOF-flavoured GeoJSON. (Pass a value of '0' or 'false' if you need to index non-WOF documents.")

	fs.Bool("enable-custom-placetypes", false, "...")
	fs.String("custom-placetypes-source", "", "...")
	fs.String("custom-placetypes", "", "...")

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
