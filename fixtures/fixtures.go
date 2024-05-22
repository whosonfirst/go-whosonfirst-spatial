package fixtures

import (
	"embed"
)

//go:embed *.geojson
var FS embed.FS
