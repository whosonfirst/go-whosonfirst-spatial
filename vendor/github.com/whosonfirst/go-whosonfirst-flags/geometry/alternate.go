package geometry

import (
	"github.com/whosonfirst/go-whosonfirst-flags"
	"github.com/whosonfirst/go-whosonfirst-uri"	
)

type AlternateGeometryFlag struct {
	flags.AlternateGeometryFlag
	is_alt bool
	label string
}

func NewAlternateGeometryFlag(uri_str string) (flags.AlternateGeometryFlag, error) {

	_, uri_args, err := uri.ParseURI(uri_str)

	if err != nil {
		return nil, err
	}

	is_alt := uri_args.IsAlternate
	alt_label := ""
	
	if  is_alt {
		
		label, err := uri_args.AltGeom.String()

		if err != nil {
			return nil, err
		}

		alt_label = label
	}
		
	// check label against go-whosonfirst-sources here?

	f := AlternateGeometryFlag{
		is_alt: is_alt,
		label: alt_label,
	}

	return &f, nil
}

func (f *AlternateGeometryFlag) MatchesAny(others ...flags.AlternateGeometryFlag) bool {

	for _, o := range others {

		if f.isEqual(o){
			return true
		}

	}

	return false
}

func (f *AlternateGeometryFlag) MatchesAll(others ...flags.AlternateGeometryFlag) bool {

	matches := 0

	for _, o := range others {

		if f.isEqual(o){
			matches += 1
		}

	}

	if matches == len(others) {
		return true
	}

	return false
}

func (f *AlternateGeometryFlag) IsAlternateGeometry() bool {
	return f.is_alt
}

func (f *AlternateGeometryFlag) Label() string {
	return f.label
}

func (f *AlternateGeometryFlag) String() string {
	return f.Label()
}

func (f *AlternateGeometryFlag) isEqual(other flags.AlternateGeometryFlag) bool {

	if f.IsAlternateGeometry() != other.IsAlternateGeometry(){
		return false
	}

	if f.Label() != other.Label(){
		return false
	}

	return true
}
