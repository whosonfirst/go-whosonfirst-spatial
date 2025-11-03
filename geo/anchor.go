package geo

// WARNING: VIBE-CODING - NONE OF THIS HAS BEEN TESTED OR VETTED YET
// Examine the code in the uploaded file which is code written in JavaScript.
// Then implement the same functionality defined in "findAnchorPoint" as a Go function that uses no external dependencies.

import (
	"math"
	"sort"
)

// ---------------------------------------------------------------------
// Basic geometry types and helpers
// ---------------------------------------------------------------------

// Point represents a 2‑D coordinate.
type Point struct{ X, Y float64 }

// Anchor contains an anchor point and its weighted distance to the edge.
type Anchor struct{ X, Y, Distance float64 }

// Bounds holds the axis‑aligned bounding box of a shape.
type Bounds struct{ xmin, xmax, ymin, ymax float64 }

func (b Bounds) Width() float64  { return b.xmax - b.xmin }
func (b Bounds) Height() float64 { return b.ymax - b.ymin }
func (b Bounds) Area() float64   { return b.Width() * b.Height() }

// ---------------------------------------------------------------------
// Geometry helpers (polygon area, centroid, bounds, etc.)
// ---------------------------------------------------------------------

// getBounds returns the bounding box of a single ring.
func getBounds(ring []Point) Bounds {
	if len(ring) == 0 {
		return Bounds{}
	}
	minX, maxX := ring[0].X, ring[0].X
	minY, maxY := ring[0].Y, ring[0].Y
	for _, p := range ring[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	return Bounds{minX, maxX, minY, maxY}
}

// signedArea returns the signed shoelace area of a ring.
func signedArea(ring []Point) float64 {
	if len(ring) < 3 {
		return 0
	}
	a := 0.0
	n := len(ring)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		a += ring[i].X*ring[j].Y - ring[j].X*ring[i].Y
	}
	return a * 0.5
}

// area returns the absolute area of a ring.
func area(ring []Point) float64 { return math.Abs(signedArea(ring)) }

// centroid returns the centroid of a ring (weighted by area).
func centroid(ring []Point) Point {
	if len(ring) < 3 {
		// Degenerate – just take the average of the points.
		x, y := 0.0, 0.0
		for _, p := range ring {
			x += p.X
			y += p.Y
		}
		n := float64(len(ring))
		return Point{x / n, y / n}
	}
	cx, cy, a := 0.0, 0.0, 0.0
	n := len(ring)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		xi, yi := ring[i].X, ring[i].Y
		xj, yj := ring[j].X, ring[j].Y
		cross := xi*yj - xj*yi
		a += cross
		cx += (xi + xj) * cross
		cy += (yi + yj) * cross
	}
	a *= 0.5
	cx /= (6 * a)
	cy /= (6 * a)
	return Point{cx, cy}
}

// getMaxPathIndex returns the index of the ring with the largest absolute area.
func getMaxPathIndex(shape [][]Point) int {
	maxIdx, maxArea := 0, 0.0
	for i, ring := range shape {
		if a := area(ring); a > maxArea {
			maxArea = a
			maxIdx = i
		}
	}
	return maxIdx
}

// ---------------------------------------------------------------------
// Point‑to‑polygon utilities
// ---------------------------------------------------------------------

// pointSegmentDistance returns the shortest distance from p to segment a‑b.
func pointSegmentDistance(p, a, b Point) float64 {
	// projection t of p onto line a‑b
	dx, dy := b.X-a.X, b.Y-a.Y
	if dx == 0 && dy == 0 {
		return math.Hypot(p.X-a.X, p.Y-a.Y)
	}
	t := ((p.X-a.X)*dx + (p.Y-a.Y)*dy) / (dx*dx + dy*dy)
	if t < 0 {
		return math.Hypot(p.X-a.X, p.Y-a.Y)
	}
	if t > 1 {
		return math.Hypot(p.X-b.X, p.Y-b.Y)
	}
	// perpendicular distance
	px := a.X + t*dx
	py := a.Y + t*dy
	return math.Hypot(p.X-px, p.Y-py)
}

// pointToShapeDistance returns the minimum distance from a point to any edge
// of the polygon (outer ring + holes).
func pointToShapeDistance(x, y float64, shape [][]Point) float64 {
	p := Point{x, y}
	minDist := math.Inf(1)
	for _, ring := range shape {
		for i := 0; i < len(ring); i++ {
			j := (i + 1) % len(ring)
			d := pointSegmentDistance(p, ring[i], ring[j])
			if d < minDist {
				minDist = d
			}
		}
	}
	return minDist
}

// ---------------------------------------------------------------------
// Point‑in‑polygon test (handles holes)
// ---------------------------------------------------------------------

// pointInRing returns true if (x,y) is inside the ring using the ray‑crossing rule.
func pointInRing(x, y float64, ring []Point) bool {
	if len(ring) == 0 {
		return false
	}
	inside := false
	for i := 0; i < len(ring); i++ {
		j := (i + 1) % len(ring)
		yi := ring[i].Y
		yj := ring[j].Y
		if ((yi > y) != (yj > y)) &&
			(x < (ring[j].X-ring[i].X)*(y-yi)/(yj-yi)+ring[i].X) {
			inside = !inside
		}
	}
	return inside
}

// testPointInPolygon returns true if the point lies in the outer ring
// and not inside any hole.
func testPointInPolygon(x, y float64, shape [][]Point) bool {
	if len(shape) == 0 {
		return false
	}
	if !pointInRing(x, y, shape[0]) {
		return false
	}
	for i := 1; i < len(shape); i++ {
		if pointInRing(x, y, shape[i]) {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------
// Ray‑intersection helpers
// ---------------------------------------------------------------------

// getRayIntersection returns the y‑coordinate where the segment (x1,y1)-(x2,y2)
// intersects the vertical line x = x0.
// If the segment is vertical or does not cross the line, returns NaN.
func getRayIntersection(x0, y0, x1, y1, x2, y2 float64) float64 {
	// We only need to know whether the segment crosses the vertical line.
	// Skip vertical segments that lie on the line (they would produce infinite
	// intersections – MapShaper ignores such cases).
	if x1 == x2 {
		return math.NaN()
	}
	if (x0 < math.Min(x1, x2)) || (x0 > math.Max(x1, x2)) {
		return math.NaN()
	}
	// linear interpolation
	yi := y1 + (y2-y1)*(x0-x1)/(x2-x1)
	return yi
}

// findRayRingIntersections returns all y‑intersections of a vertical ray
// at x = x0 with the ring.
func findRayRingIntersections(x0, y0 float64, ring []Point) []float64 {
	var yints []float64
	for i := 0; i < len(ring); i++ {
		j := (i + 1) % len(ring)
		yi := getRayIntersection(x0, y0, ring[i].X, ring[i].Y, ring[j].X, ring[j].Y)
		if !math.IsNaN(yi) {
			yints = append(yints, yi)
		}
	}
	// If the ray touches an odd number of edges (e.g. grazes a vertex),
	// ignore the result – we need an even count to form interior segments.
	if len(yints)%2 == 1 {
		return nil
	}
	return yints
}

// findRayShapeIntersections aggregates all ring intersections for the shape.
func findRayShapeIntersections(x0, y0 float64, shape [][]Point) []float64 {
	var yints []float64
	for _, ring := range shape {
		yints = append(yints, findRayRingIntersections(x0, y0, ring)...)
	}
	return yints
}

// ---------------------------------------------------------------------
// Candidate generation
// ---------------------------------------------------------------------

// candidate holds a point and the half‑segment length that produced it.
type candidate struct {
	X, Y, Interval float64
}

// findHitCandidates returns the mid‑points of all interior segments
// produced by intersecting the vertical line x = x0 with the polygon.
func findHitCandidates(x0 float64, shape [][]Point, globalMinY float64) []candidate {
	yints := findRayShapeIntersections(x0, globalMinY, shape)
	if len(yints) == 0 {
		return nil
	}
	sort.Float64s(yints)
	var cands []candidate
	for i := 0; i+1 < len(yints); i += 2 {
		y1, y2 := yints[i], yints[i+1]
		interval := (y2 - y1) / 2
		if interval > 0 {
			cands = append(cands, candidate{
				X:       x0,
				Y:       (y1 + y2) / 2,
				Interval: interval,
			})
		}
	}
	return cands
}

// findAnchorPointCandidates evaluates a set of x‑values and collects all
// hit candidates.
func findAnchorPointCandidates(shape [][]Point, xs []float64) []candidate {
	// Need a global ymin to start the ray far below the polygon.
	// Compute it once here.
	globalMinY := math.Inf(1)
	for _, ring := range shape {
		for _, p := range ring {
			if p.Y < globalMinY {
				globalMinY = p.Y
			}
		}
	}
	globalMinY -= 1

	var all []candidate
	for _, x := range xs {
		all = append(all, findHitCandidates(x, shape, globalMinY)...)
	}
	return all
}

// ---------------------------------------------------------------------
// Weighting function
// ---------------------------------------------------------------------

// getPointWeightingFunction returns a closure that gives a weight
// (0.75–1) based on distance from the centroid.
func getPointWeightingFunction(centroid Point, bounds Bounds) func(x, y float64) float64 {
	refDist := math.Max(bounds.Width(), bounds.Height()) / 2
	return func(x, y float64) float64 {
		offset := math.Hypot(x-centroid.X, y-centroid.Y)
		return 1 - math.Min(0.6*offset/refDist, 0.25)
	}
}

// ---------------------------------------------------------------------
// Anchor point search
// ---------------------------------------------------------------------

// getInnerTics returns n evenly spaced points between min and max
// (excluding the endpoints).
func getInnerTics(min, max float64, steps int) []float64 {
	if steps <= 0 {
		return nil
	}
	step := (max - min) / float64(steps+1)
	tics := make([]float64, steps)
	for i := 1; i <= steps; i++ {
		tics[i-1] = min + step*float64(i)
	}
	return tics
}

// getAdjustedPoint moves the point vertically to maximise its weighted distance
// to the polygon edge.
func getAdjustedPoint(x, y float64, shape [][]Point, vstep float64, weight func(x, y float64) float64) Anchor {
	a := Anchor{
		X:       x,
		Y:       y,
		Distance: pointToShapeDistance(x, y, shape) * weight(x, y),
	}
	// scan up
	anchor := a
	anchor = scanForBetterPoint(anchor, shape, vstep, weight)
	// scan down
	anchor = scanForBetterPoint(anchor, shape, -vstep, weight)
	return anchor
}

// scanForBetterPoint iteratively moves the point along the vertical
// axis (up if step>0, down if step<0) until the weighted distance
// stops improving.
func scanForBetterPoint(p Anchor, shape [][]Point, vstep float64, weight func(x, y float64) float64) Anchor {
	x, y := p.X, p.Y
	bestD := p.Distance
	for {
		y += vstep
		d := pointToShapeDistance(x, y, shape) * weight(x, y)
		// allow a small tolerance for local minima
		if d > bestD*0.90 && testPointInPolygon(x, y, shape) {
			if d > bestD {
				bestD = d
				p.Y = y
			}
		} else {
			break
		}
	}
	p.Distance = bestD
	return p
}

// probeForBestAnchorPoint evaluates a set of candidate points and returns
// the one with the largest weighted distance to the polygon edge.
func probeForBestAnchorPoint(shape [][]Point, lbound, rbound float64, htics int,
	weight func(x, y float64) float64) *Anchor {

	tics := getInnerTics(lbound, rbound, htics)
	interval := (rbound - lbound) / float64(htics)

	cands := findAnchorPointCandidates(shape, tics)

	// weight each candidate by the half‑segment length and the distance weighting
	for i := range cands {
		cands[i].Interval *= weight(cands[i].X, cands[i].Y)
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].Interval > cands[j].Interval })

	var best *Anchor
	for _, cand := range cands {
		if best != nil && best.Distance > cand.Interval {
			// remaining candidates can’t beat best
			break
		}
		adj := getAdjustedPoint(cand.X, cand.Y, shape, interval, weight)
		if best == nil || adj.Distance > best.Distance {
			best = &adj
		}
	}
	return best
}

// findAnchorPoint2 contains the core logic that picks a suitable anchor point.
func findAnchorPoint2(shape [][]Point) *Anchor {
	if len(shape) == 0 {
		return nil
	}
	maxIdx := getMaxPathIndex(shape)
	maxPath := shape[maxIdx]
	bounds := getBounds(maxPath)
	centroid := centroid(maxPath)
	weight := getPointWeightingFunction(centroid, bounds)

	area := area(maxPath)

	var htics int
	var focus float64
	if len(shape) == 1 && area*1.2 > bounds.Area() {
		htics = 5
		focus = 0.2
	} else if len(shape) == 1 && area*1.7 > bounds.Area() {
		htics = 7
		focus = 0.4
	} else {
		htics = 11
		focus = 0.5
	}

	hRange := bounds.Width() * focus
	lbound := centroid.X - hRange/2
	rbound := lbound + hRange
	hstep := hRange / float64(htics)

	best := probeForBestAnchorPoint(shape, lbound, rbound, htics, weight)
	if best == nil {
		// fallback to centroid
		return &Anchor{
			X:       centroid.X,
			Y:       centroid.Y,
			Distance: pointToShapeDistance(centroid.X, centroid.Y, shape) * weight(centroid.X, centroid.Y),
		}
	}

	// Look for an even better fit in a small window around the best point.
	p2 := probeForBestAnchorPoint(shape, best.X-hstep/2, best.X+hstep/2, 2, weight)
	if p2 != nil && p2.Distance > best.Distance {
		best = p2
	}
	return best
}

// ---------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------

// FindAnchorPoint returns an anchor point inside the polygon shape that is
// maximally distant from the polygon edge, weighted by proximity to the centroid.
// The shape is represented as [][]Point – the first slice is the outer ring,
// any subsequent slices are holes.
//
// The function is a direct translation of the MapShaper JavaScript routine
// and has no external dependencies.
func FindAnchorPoint(shape [][]Point) *Anchor {
	return findAnchorPoint2(shape)
}
