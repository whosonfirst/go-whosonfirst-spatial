package main

import (
	"context"
	"log"

	"github.com/whosonfirst/go-whosonfirst-spatial/app/pip"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := pip.Run(ctx, logger)

	if err != nil {
		logger.Fatal(err)
	}
}
