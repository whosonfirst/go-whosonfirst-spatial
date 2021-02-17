package flags

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-index/v2/emitter"
	"sort"
	"strings"
)

func AppendIndexingFlags(fs *flag.FlagSet) error {

	modes := emitter.Schemes()
	sort.Strings(modes)

	valid_modes := strings.Join(modes, ", ")
	desc_modes := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-index/v2/emitter URI. Supported schemes are: %s.", valid_modes)

	fs.String("emitter-uri", "repo://", desc_modes)

	return nil
}

func ValidateIndexingFlags(fs *flag.FlagSet) error {

	_, err := StringVar(fs, "mode")

	if err != nil {
		return err
	}

	return nil
}
