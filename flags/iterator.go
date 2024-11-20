package flags

import (
	"fmt"
	"strings"
)

const SEP_FRAGMENT string = "#"
const SEP_PIPE string = "|"
const SEP_SPACE string = " "

type MultiIteratorURIFlag []*IteratorURIFlag

func (fl *MultiIteratorURIFlag) Set(value string) error {

	iter_fl := new(IteratorURIFlag)

	err := iter_fl.Set(value)

	if err != nil {
		return err
	}

	*fl = append(*fl, iter_fl)
	return nil
}

func (fl *MultiIteratorURIFlag) Key() string {
	return ""
}

func (fl *MultiIteratorURIFlag) Value() interface{} {
	return fl
}

func (fl *MultiIteratorURIFlag) String() string {

	str_flags := make([]string, len(*fl))

	for i, iter_fl := range *fl {
		str_flags[i] = iter_fl.String()
	}

	return strings.Join(str_flags, SEP_SPACE)
}

func (fl *MultiIteratorURIFlag) AsMap() map[string][]string {

	iter_map := make(map[string][]string)

	for _, iter_fl := range *fl {

		iter_uri := iter_fl.Key()
		iter_sources := iter_fl.Value().([]string)

		iter_map[iter_uri] = iter_sources
	}

	return iter_map
}

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
