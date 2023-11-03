package database

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"github.com/whosonfirst/go-whosonfirst-feature/geometry"
)

func IndexDatabaseWithFS(ctx context.Context, db SpatialDatabase, index_fs fs.FS) error {

	walk_func := func(path string, d fs.DirEntry, err error) error {

		if d.IsDir() {
			return nil
		}

		r, err := index_fs.Open(path)

		if err != nil {
			return fmt.Errorf("Failed to open %s for reading, %w", path, err)
		}

		defer r.Close()

		body, err := io.ReadAll(r)

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", path, err)
		}

		geom_type, err := geometry.Type(body)

		if err != nil {
			return fmt.Errorf("Failed to derive geometry type for %s, %w", path, err)
		}

		switch geom_type {
		case "Polygon", "MultiPolygon":
			return db.IndexFeature(ctx, body)
		default:
			return nil
		}
		return nil
	}

	return fs.WalkDir(index_fs, ".", walk_func)
}
