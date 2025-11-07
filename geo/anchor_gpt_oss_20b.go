package geo

/*

Examine these two files. The first is JavaScript code from the Mapshaper project to derive the most suitable centroid from a Polygon or MultiPolygon. The second file tries to implement this same logic in Go. Analyze the latter and suggest an improved version to better implement the latter.

mapshaper-anchor-points.mjs
anchor_qwen3_coder_30b.go

*/

import (
	"math"
	"sort"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/simplify"
)

// FindAnchorPoint returns a point that lies inside a polygon or multipolygon
// and is as far from the boundary as possible.
func FindAnchorPoint(g orb.Geometry) (orb.Point, error) {
	switch geom := g.(type) {
	case orb.Polygon:
		return findAnchorPointInPolygon(geom)
	case orb.MultiPolygon:
		return findAnchorPointInMultiPolygon(geom)
	default:
		return orb.Point{}, nil // non‑polygon geometries get an empty point
	}
}

// ---------------------------------------------------------------------------
// Polygon handling
// ---------------------------------------------------------------------------

func findAnchorPointInPolygon(poly orb.Polygon) (orb.Point, error) {
	if len(poly) == 0 {
		return orb.Point{}, nil
	}

	outer := getLargestRing(poly)
	if len(outer) < 3 {
		return orb.Point{}, nil
	}

	// Bounds & area
	bounds := outer.Bound()
	if bounds.Min[0] == bounds.Max[0] && bounds.Min[1] == bounds.Max[1] {
		return orb.Point{}, nil
	}
	area := planar.Area(outer)

	// Simplify the polygon
	thresh := math.Sqrt(area) * 0.01
	simplified := simplify.DouglasPeucker(thresh).Ring(outer)
	if len(simplified) < 3 {
		// Very thin polygons – fall back to centroid
		c, _ := planar.CentroidArea(orb.Polygon{simplified})
		return c, nil
	}

	return findAnchorPoint2(simplified, bounds)
}

// ---------------------------------------------------------------------------
// MultiPolygon handling
// ---------------------------------------------------------------------------

func findAnchorPointInMultiPolygon(multi orb.MultiPolygon) (orb.Point, error) {
	if len(multi) == 0 {
		return orb.Point{}, nil
	}

	var bestPoly orb.Polygon
	var bestArea float64

	for _, poly := range multi {
		if len(poly) == 0 {
			continue
		}
		outer := getLargestRing(poly)
		if len(outer) < 3 {
			continue
		}
		area := planar.Area(outer)
		if area > bestArea {
			bestArea = area
			bestPoly = poly
		}
	}

	if len(bestPoly) > 0 {
		return findAnchorPointInPolygon(bestPoly)
	}

	// Fallback: centroid of the first polygon that has an outer ring
	for _, poly := range multi {
		if len(poly) > 0 {
			outer := getLargestRing(poly)
			if len(outer) >= 3 {
				c, _ := planar.CentroidArea(orb.Polygon{outer})
				return c, nil
			}
		}
	}
	return orb.Point{}, nil
}

// ---------------------------------------------------------------------------
// Core anchor‑point logic
// ---------------------------------------------------------------------------

func findAnchorPoint2(path orb.Ring, bounds orb.Bound) (orb.Point, error) {
	centroid, _ := planar.CentroidArea(orb.Polygon{path})
	weight := getPointWeightingFunction(centroid, bounds)
	area := planar.Area(path)

	// Adapt number of horizontal tics depending on how “simple” the shape is
	var htics int
	var focus float64
	boundArea := getBoundArea(bounds)

	switch {
	case area*1.2 > boundArea:
		htics = 5
		focus = 0.2
	case area*1.7 > boundArea:
		htics = 7
		focus = 0.4
	default:
		htics = 11
		focus = 0.5
	}

	hrange := getBoundWidth(bounds) * focus
	lbound := centroid[0] - hrange/2
	rbound := lbound + hrange
	hstep := hrange / float64(htics)

	// First pass – a wide search
	best, _ := probeForBestAnchorPoint(path, lbound, rbound, htics, weight, hstep)
	if best.Distance == 0 {
		// Failed – fall back to centroid
		return centroid, nil
	}

	// Second pass – a fine search around the best point
	best2, _ := probeForBestAnchorPoint(path, best.X-hstep/2, best.X+hstep/2, 2, weight, hstep/2)
	if best2.Distance > best.Distance {
		best = best2
	}

	return orb.Point{best.X, best.Y}, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func getLargestRing(poly orb.Polygon) orb.Ring {
	if len(poly) == 0 {
		return orb.Ring{}
	}
	largest := poly[0]
	for _, r := range poly {
		if len(r) > len(largest) {
			largest = r
		}
	}
	return largest
}

func getPointWeightingFunction(centroid orb.Point, bounds orb.Bound) func(orb.Point) float64 {
	refDist := math.Max(getBoundWidth(bounds), getBoundHeight(bounds)) / 2
	return func(p orb.Point) float64 {
		offset := math.Hypot(centroid[0]-p[0], centroid[1]-p[1])
		weight := 1 - math.Min(0.6*offset/refDist, 0.25)
		return math.Max(0, weight)
	}
}

func getInnerTics(min, max float64, steps int) []float64 {
	if steps <= 0 {
		return nil
	}
	rangeVal := max - min
	step := rangeVal / float64(steps+1)
	tics := make([]float64, 0, steps)
	for i := 1; i <= steps; i++ {
		tics = append(tics, min+step*float64(i))
	}
	return tics
}

func findAnchorPointCandidates(path orb.Ring, tics []float64) []Candidate {
	var candidates []Candidate
	for _, x := range tics {
		candidates = append(candidates, findHitCandidates(x, path)...)
	}
	return candidates
}

func findHitCandidates(x float64, path orb.Ring) []Candidate {
	if len(path) < 2 {
		return nil
	}

	// 1. Find all vertical intersections
	var yints []float64
	for i := 0; i < len(path)-1; i++ {
		p1, p2 := path[i], path[i+1]
		if intersectsVertical(p1, p2, x) {
			yints = append(yints, verticalIntersection(p1, p2, x))
		}
	}
	// 2. Close the ring
	if len(path) > 2 {
		p1, p2 := path[len(path)-1], path[0]
		if intersectsVertical(p1, p2, x) {
			yints = append(yints, verticalIntersection(p1, p2, x))
		}
	}
	if len(yints) == 0 {
		return nil
	}
	sort.Float64s(yints)

	// 3. Pair intersections into segments and take midpoints
	var cands []Candidate
	for i := 0; i < len(yints)-1; i += 2 {
		y1, y2 := yints[i], yints[i+1]
		if y2 <= y1 {
			continue
		}
		interval := (y2 - y1) / 2 // half vertical segment
		cands = append(cands, Candidate{
			X:        x,
			Y:        (y1 + y2) / 2,
			Interval: interval,
		})
	}
	return cands
}

func intersectsVertical(p1, p2 orb.Point, x float64) bool {
	return (p1[0] <= x && x <= p2[0]) || (p2[0] <= x && x <= p1[0])
}

func verticalIntersection(p1, p2 orb.Point, x float64) float64 {
	// segment is not vertical (otherwise it would be a touch)
	return p1[1] + (p2[1]-p1[1])*(x-p1[0])/(p2[0]-p1[0])
}

type Candidate struct {
	X, Y     float64
	Interval float64 // half vertical segment length
	Distance float64 // weighted distance to polygon boundary
}

func probeForBestAnchorPoint(path orb.Ring, lbound, rbound float64, htics int, weight func(orb.Point) float64, vstep float64) (Candidate, error) {
	tics := getInnerTics(lbound, rbound, htics)

	candidates := findAnchorPointCandidates(path, tics)

	// Weight each candidate by the interval * centroid weighting
	for i := range candidates {
		candidates[i].Interval *= weight(orb.Point{candidates[i].X, candidates[i].Y})
	}

	// Sort by weighted interval descending
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Interval > candidates[j].Interval
	})

	var best Candidate
	for _, cand := range candidates {
		// Early break if the weighted interval of the next candidate
		// cannot beat the current best distance
		if best.Distance > 0 && best.Distance > cand.Interval {
			break
		}

		adj := getAdjustedPoint(cand.X, cand.Y, path, vstep, weight)
		if adj.Distance > best.Distance {
			best = adj
		}
	}

	return best, nil
}

func getAdjustedPoint(x, y float64, path orb.Ring, vstep float64, weight func(orb.Point) float64) Candidate {
	pt := orb.Point{x, y}
	dist := planar.DistanceFrom(path, pt) * weight(pt)
	cand := Candidate{X: x, Y: y, Distance: dist}

	// Scan upward
	scanForBetterPoint(&cand, path, vstep, weight, 1)
	// Scan downward
	scanForBetterPoint(&cand, path, vstep, weight, -1)

	return cand
}

func scanForBetterPoint(c *Candidate, path orb.Ring, step float64, weight func(orb.Point) float64, dir float64) {
	maxDist := c.Distance
	x := c.X
	y := c.Y

	for {
		y += dir * step
		pt := orb.Point{x, y}
		if !planar.RingContains(path, pt) {
			break
		}
		dist := planar.DistanceFrom(path, pt) * weight(pt)
		// Stop when the distance falls below 90 % of the best found
		if dist > maxDist*0.9 {
			if dist > maxDist {
				maxDist = dist
				c.Y = y
				c.Distance = dist
			}
		} else {
			break
		}
	}
}

// ---------------------------------------------------------------------------
// Bound helpers
// ---------------------------------------------------------------------------

func getBoundArea(b orb.Bound) float64   { return getBoundWidth(b) * getBoundHeight(b) }
func getBoundWidth(b orb.Bound) float64  { return b.Max[0] - b.Min[0] }
func getBoundHeight(b orb.Bound) float64 { return b.Max[1] - b.Min[1] }

// ---------------------------------------------------------------------------
// Public wrapper for backward compatibility
// ---------------------------------------------------------------------------

func FindAnchorPointWrapper(g orb.Geometry) (orb.Point, error) {
	return FindAnchorPoint(g)
}
