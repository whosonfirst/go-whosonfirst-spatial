package flags

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-iterate/emitter"
	"sort"
	"strings"
)

func AppendIndexingFlags(fs *flag.FlagSet) error {

	modes := emitter.Schemes()
	sort.Strings(modes)

	valid_modes := strings.Join(modes, ", ")
	desc_modes := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterate/emitter URI. Supported schemes are: %s.", valid_modes)

	fs.String("iterator-uri", "repo://", desc_modes)

	return nil
}

func ValidateIndexingFlags(fs *flag.FlagSet) error {

	_, err := StringVar(fs, "iterator-uri")

	if err != nil {
		return err
	}

	return nil
}
