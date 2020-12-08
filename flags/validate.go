package flags

import (
	"flag"
	"errors"
)

func ValidateCommonFlags(fs *flag.FlagSet) error {

	_, err := StringVar(fs, "mode")

	if err != nil {
		return err
	}

	_, err = StringVar(fs, "spatial-database-uri")

	if err != nil {
		return err
	}

	enable_properties, err := BoolVar(fs, "enable-properties")

	if err != nil {
		return err
	}

	if enable_properties {

		properties_reader_uri, err := StringVar(fs, "properties-reader-uri")

		if err != nil {
			return err
		}

		if properties_reader_uri == "" {
			return errors.New("Invalid or missing -properties-reader-uri flag")
		}
	}

	return nil
}
