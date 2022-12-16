# go-whosonfirst-spatial

Go package defining interfaces for Who's On First specific spatial operations.

## Documentation

Documentation, particularly proper Go documentation, is incomplete at this time.

## Motivation

The goal of the `go-whosonfirst-spatial` package is to de-couple the various components that make up the [go-whosonfirst-pip-v2](https://github.com/whosonfirst/go-whosonfirst-pip-v2) package – indexing, storage, querying and serving – in to separate packages in order to allow for more flexibility.

It is the "base" package that defines provider-agnostic, but WOF-specific, interfaces for a limited set of spatial queries and reading properties.

These interfaces are then implemented in full or in part by provider-specific classes. For example, an in-memory RTree index or a SQLite database or even a Protomaps database:

* https://github.com/whosonfirst/go-whosonfirst-spatial-rtree
* https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite
* https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles

Building on that there are equivalent base packages for "server" implementations, like:

* https://github.com/whosonfirst/go-whosonfirst-spatial-www
* https://github.com/whosonfirst/go-whosonfirst-spatial-grpc

The idea is that all of these pieces can be _easily_ combined in to purpose-fit applications.  As a practical matter it's mostly about trying to identify and package the common pieces in to as few lines of code as possible so that they might be combined with an application-specific `import` statement. For example:

```
import (
         _ "github.com/whosonfirst/go-whosonfirst-spatial-MY-SPECIFIC-REQUIREMENTS"
)
```

Here is a concrete example, implementing a point-in-polygon service over HTTP using a SQLite backend:

* https://github.com/whosonfirst/go-whosonfirst-spatial-www/blob/main/application/server
* https://github.com/whosonfirst/go-whosonfirst-spatial-www-sqlite/blob/main/cmd/server/main.go

It is part of the overall goal of:

* Staying out people's database or delivery choices (or needs)
* Supporting as many databases (and delivery (and indexing) choices) as possible
* Not making database `B` a dependency (in the Go code) in order to use database `A`, as in not bundling everything in a single mono-repo that becomes bigger and has more requirements over time.

Importantly this package does not implement any actual spatial functionality. It defines the interfaces that are implemented by other packages which allows code to function without the need to consider the underlying mechanics of how spatial operations are being performed.

## Concepts

### SpatialIndex


### SpatialDatabase

Any system that can store and query for one or more Who's On First record, implementing the `database.SpatialDatabase` interface.

### Properties Reader

_Please write me_

### Filters

_Please write me_

### Standard Places Response (SPR)

_Please write me_

## Implementations

* https://github.com/whosonfirst/go-whosonfirst-spatial-rtree
* https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite
* https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles

## Servers and clients

### WWW

* https://github.com/whosonfirst/go-whosonfirst-spatial-www
* https://github.com/whosonfirst/go-whosonfirst-spatial-www-sqlite
* https://github.com/whosonfirst/go-whosonfirst-spatial-www-pmtiles

### gRPC

* https://github.com/whosonfirst/go-whosonfirst-spatial-grpc
* https://github.com/whosonfirst/go-whosonfirst-spatial-grpc-sqlite
* https://github.com/whosonfirst/go-whosonfirst-spatial-grpc-pmtiles

## Services and Operations

* https://github.com/whosonfirst/go-whosonfirst-spatial-pip
* https://github.com/whosonfirst/go-whosonfirst-spatial-hierarchy