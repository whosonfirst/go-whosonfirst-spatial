package flags

import (
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/emitter"
	"sort"
	"strings"
)

// AppendIndexingFlag will append indexing (whosonfirst/go-whosonfirst-iterate/v2) related flags to 'fs'.
func AppendIndexingFlags(fs *flag.FlagSet) error {

	modes := emitter.Schemes()
	sort.Strings(modes)

	valid_modes := strings.Join(modes, ", ")
	desc_modes := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterate/v2 URI. Supported schemes are: %s.", valid_modes)

	fs.String(IteratorURIFlag, "repo://", desc_modes)

	return nil
}

// ValidateIndexingFlags will ensure that all indexing (whosonfirst/go-whosonfirst-iterate/v2) related flags have been assigned to 'fs'.
func ValidateIndexingFlags(fs *flag.FlagSet) error {

	_, err := lookup.StringVar(fs, IteratorURIFlag)

	if err != nil {
		return fmt.Errorf("Failed to lookup %s flag, %w", IteratorURIFlag, err)
	}

	return nil
}
