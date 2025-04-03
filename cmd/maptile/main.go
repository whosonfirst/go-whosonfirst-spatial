package main

import (
	"flag"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/maptile"
)

func main() {

	var latitude float64
	var longitude float64
	var zoom int

	flag.Float64Var(&latitude, "latitude", 0.0, "")
	flag.Float64Var(&longitude, "longitude", 0.0, "")
	flag.IntVar(&zoom, "zoom", 0, "")

	flag.Parse()

	pt := orb.Point([2]float64{longitude, latitude})
	z := maptile.Zoom(uint32(zoom))
	t := maptile.At(pt, z)

	fmt.Printf("%d/%d/%d", t.Z, t.X, t.Y)
}
