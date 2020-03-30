package database

import (
	"context"
)

type ExtrasDatabase interface {
	Close(context.Context) error
}
