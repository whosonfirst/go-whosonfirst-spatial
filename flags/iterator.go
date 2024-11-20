package flags

import (
	"fmt"
	"strings"
)

const SEP_FRAGMENT string = "#"
const SEP_PIPE string = "|"

type IteratorURIFlag struct {
	iter_uri     string
	iter_sources []string
}

func (fl *IteratorURIFlag) Key() string {
	return fl.iter_uri
}

func (fl *IteratorURIFlag) Value() interface{} {
	return fl.iter_sources
}

func (fl *IteratorURIFlag) String() string {

	str_sources := strings.Join(fl.iter_sources, SEP_PIPE)

	parts := []string{
		fl.iter_uri,
		str_sources,
	}

	return strings.Join(parts, SEP_FRAGMENT)
}

func (fl *IteratorURIFlag) Set(value string) error {

	parts := strings.Split(value, SEP_FRAGMENT)

	if len(parts) != 2 {
		return fmt.Errorf("Invalid iterator URI")
	}

	iter_uri := parts[0]
	iter_sources := strings.Split(parts[1], SEP_PIPE)

	if len(iter_sources) == 0 {
		return fmt.Errorf("Iterator URI missing sources")
	}

	fl.iter_uri = iter_uri
	fl.iter_sources = iter_sources

	return nil
}
