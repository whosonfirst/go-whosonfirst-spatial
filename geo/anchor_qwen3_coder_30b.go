//go:build qwen

package geo

// Correct:
// mapshaper -i /usr/local/data/sfomuseum-data-whosonfirst/data/856/886/37/85688637.geojson -points inner -o -
// "geometry": {"type":"Point","coordinates":[-119.58616773283674,36.53154412128562]}

// Wrong:
// go run cmd/anchor/main.go /usr/local/data/sfomuseum-data-whosonfirst/data/856/886/37/85688637.geojson
// "geometry": {"type":"Point","coordinates":[-119.0413633896104,33.48935415584416]}

import (
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/simplify"
)

// FindAnchorPoint finds a point inside a polygon or multipolygon that is away from the polygon edge
func FindAnchorPoint(g orb.Geometry) (orb.Point, error) {
	switch geom := g.(type) {
	case orb.Polygon:
		return findAnchorPointInPolygon(geom)
	case orb.MultiPolygon:
		return findAnchorPointInMultiPolygon(geom)
	default:
		return orb.Point{}, nil // Return empty point for non-polygon geometries
	}
}

// findAnchorPointInPolygon finds anchor point in a single polygon
func findAnchorPointInPolygon(poly orb.Polygon) (orb.Point, error) {
	if len(poly) == 0 {
		return orb.Point{}, nil
	}

	// Get the largest ring (outer ring)
	maxPath := getLargestRing(poly)
	if len(maxPath) < 3 {
		return orb.Point{}, nil
	}

	// Calculate bounds
	// bounds := getBounds(maxPath)
	bounds := maxPath.Bound()

	if bounds.Min[0] == bounds.Max[0] && bounds.Min[1] == bounds.Max[1] {
		return orb.Point{}, nil
	}

	// Simplify polygon
	simplified := simplifyPolygon(maxPath, bounds)
	if len(simplified) < 3 {
		return orb.Point{}, nil
	}

	// Find anchor point in simplified polygon
	return findAnchorPoint2(simplified, bounds)
}

// findAnchorPointInMultiPolygon finds anchor point in a multipolygon
func findAnchorPointInMultiPolygon(multiPoly orb.MultiPolygon) (orb.Point, error) {
	if len(multiPoly) == 0 {
		return orb.Point{}, nil
	}

	// Find the largest polygon in the multipolygon
	var largestPoly orb.Polygon
	var maxArea float64

	for _, poly := range multiPoly {
		if len(poly) == 0 {
			continue
		}

		// Get the largest ring (outer ring)
		maxPath := getLargestRing(poly)
		if len(maxPath) < 3 {
			continue
		}

		// Calculate area of this polygon
		area := getPlanarPathArea(maxPath)
		if area > maxArea {
			maxArea = area
			largestPoly = poly
		}
	}

	// If we found a valid polygon, process it
	if len(largestPoly) > 0 {
		return findAnchorPointInPolygon(largestPoly)
	}

	// Fallback: return centroid of first polygon if no valid polygon found
	if len(multiPoly) > 0 && len(multiPoly[0]) > 0 {
		maxPath := getLargestRing(multiPoly[0])
		if len(maxPath) >= 3 {
			centroid := getCentroid(maxPath)
			return centroid, nil
		}
	}

	return orb.Point{}, nil
}

// findAnchorPoint2 is the main implementation for finding anchor point
func findAnchorPoint2(path orb.Ring, bounds orb.Bound) (orb.Point, error) {
	centroid := getCentroid(path)
	weight := getPointWeightingFunction(centroid, bounds)
	area := getPlanarPathArea(path)

	var htics int
	var focus float64

	// Limit test area if shape is simple and squarish
	if area*1.2 > getBoundArea(bounds) {
		htics = 5
		focus = 0.2
	} else if area*1.7 > getBoundArea(bounds) {
		htics = 7
		focus = 0.4
	} else {
		htics = 11
		focus = 0.5
	}

	hrange := getBoundWidth(bounds) * focus
	lbound := centroid[0] - hrange/2
	rbound := lbound + hrange
	hstep := hrange / float64(htics)

	// Find a best-fit point
	p, err := probeForBestAnchorPoint(path, lbound, rbound, htics, weight)
	if err != nil || p == (orb.Point{}) {
		// Fallback to centroid
		return centroid, nil
	}

	// Look for even better fit close to best-fit point
	p2, err := probeForBestAnchorPoint(path, p[0]-hstep/2, p[0]+hstep/2, 2, weight)
	if err != nil {
		return p, nil
	}

	if p2 != (orb.Point{}) && p2[1] > p[1] {
		return p2, nil
	}

	return p, nil
}

// getLargestRing returns the largest ring from a polygon
func getLargestRing(poly orb.Polygon) orb.Ring {
	if len(poly) == 0 {
		return orb.Ring{}
	}

	largest := poly[0]
	for _, ring := range poly {
		if len(ring) > len(largest) {
			largest = ring
		}
	}

	return largest
}

// getBounds returns the bounding box of a ring
/*
func getBounds(ring orb.Ring) orb.Bound {

	if len(ring) == 0 {
		return orb.Bound{}
	}

	bounds := orb.Bound{
		Min: ring[0],
		Max: ring[0],
	}

	for _, p := range ring {
		if p[0] < bounds.Min[0] {
			bounds.Min[0] = p[0]
		}
		if p[1] < bounds.Min[1] {
			bounds.Min[1] = p[1]
		}
		if p[0] > bounds.Max[0] {
			bounds.Max[0] = p[0]
		}
		if p[1] > bounds.Max[1] {
			bounds.Max[1] = p[1]
		}
	}

	return bounds
}
*/

// getCentroid calculates the centroid of a ring
func getCentroid(ring orb.Ring) orb.Point {
	if len(ring) == 0 {
		return orb.Point{}
	}

	var xSum, ySum float64
	for _, p := range ring {
		xSum += p[0]
		ySum += p[1]
	}

	return orb.Point{
		xSum / float64(len(ring)),
		ySum / float64(len(ring)),
	}
}

// getPointWeightingFunction returns a function that weights points based on distance from centroid
func getPointWeightingFunction(centroid orb.Point, bounds orb.Bound) func(orb.Point) float64 {
	referenceDist := math.Max(getBoundWidth(bounds), getBoundHeight(bounds)) / 2
	return func(p orb.Point) float64 {
		offset := math.Hypot(centroid[0]-p[0], centroid[1]-p[1])
		weight := 1 - math.Min(0.6*offset/referenceDist, 0.25)
		return math.Max(0, weight) // Ensure non-negative weight
	}
}

// getPlanarPathArea calculates the area of a path
func getPlanarPathArea(path orb.Ring) float64 {
	if len(path) < 3 {
		return 0
	}

	return planar.Area(path)
}

// simplifyPolygon simplifies a polygon using a threshold
func simplifyPolygon(path orb.Ring, bounds orb.Bound) orb.Ring {
	// Simplification threshold based on area
	if getBoundArea(bounds) <= 0 {
		return path
	}

	return simplify.DouglasPeucker(0.0).Ring(path.Clone())
}

// probeForBestAnchorPoint finds the best anchor point by probing candidates
func probeForBestAnchorPoint(path orb.Ring, lbound, rbound float64, htics int, weight func(orb.Point) float64) (orb.Point, error) {
	tics := getInnerTics(lbound, rbound, htics)

	// Find candidates
	candidates := findAnchorPointCandidates(path, tics)

	// Sort candidates
	for i := range candidates {
		candidates[i].Interval *= weight(orb.Point{candidates[i].X, candidates[i].Y})
	}

	// Sort by interval (descending)
	for i := 0; i < len(candidates)-1; i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[i].Interval < candidates[j].Interval {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}

	var bestPoint orb.Point
	var bestDistance float64

	for _, cand := range candidates {
		// Get adjusted point
		adjusted := getAdjustedPoint(cand.X, cand.Y, path, 1.0, weight)

		if adjusted.Distance > bestDistance {
			bestDistance = adjusted.Distance
			bestPoint = orb.Point{adjusted.X, adjusted.Y}
		}
	}

	return bestPoint, nil
}

// findAnchorPointCandidates finds candidate points for anchor point
func findAnchorPointCandidates(path orb.Ring, tics []float64) []Candidate {
	var candidates []Candidate

	for _, x := range tics {
		candPoints := findHitCandidates(x, path)
		candidates = append(candidates, candPoints...)
	}

	return candidates
}

// findHitCandidates returns points at midpoints of line segments formed by intersection
func findHitCandidates(x float64, path orb.Ring) []Candidate {

	if len(path) < 2 {
		return []Candidate{}
	}

	var candidates []Candidate

	// Iterate through consecutive pairs of points in the ring
	for i := 0; i < len(path)-1; i++ {
		point1 := path[i]
		point2 := path[i+1]

		// Check if the x-coordinate falls within the range of this segment
		minX := math.Min(point1.X(), point2.X())
		maxX := math.Max(point1.X(), point2.X())

		if x >= minX && x <= maxX {
			// Calculate the y-coordinate at x using linear interpolation
			// Only proceed if we have a valid segment (not vertical)
			if point2.X() != point1.X() {
				// Linear interpolation formula: y = y1 + (y2-y1)/(x2-x1) * (x-x1)
				y := point1.Y() + (point2.Y()-point1.Y())/(point2.X()-point1.X())*(x-point1.X())

				// Calculate midpoint of the segment
				midX := (point1.X() + point2.X()) / 2.0
				midY := (point1.Y() + point2.Y()) / 2.0

				// Calculate distance from midpoint to the given x coordinate
				distance := math.Sqrt(math.Pow(midX-x, 2) + math.Pow(midY-y, 2))

				candidates = append(candidates, Candidate{
					X:        midX,
					Y:        midY,
					Interval: math.Abs(point2.X() - point1.X()),
					Distance: distance,
				})
			}
		}
	}

	// Handle the case where the path is a closed ring (last point connects to first)
	if len(path) > 2 {
		point1 := path[len(path)-1]
		point2 := path[0]

		minX := math.Min(point1.X(), point2.X())
		maxX := math.Max(point1.X(), point2.X())

		if x >= minX && x <= maxX {
			if point2.X() != point1.X() {
				y := point1.Y() + (point2.Y()-point1.Y())/(point2.X()-point1.X())*(x-point1.X())

				midX := (point1.X() + point2.X()) / 2.0
				midY := (point1.Y() + point2.Y()) / 2.0

				distance := math.Sqrt(math.Pow(midX-x, 2) + math.Pow(midY-y, 2))

				candidates = append(candidates, Candidate{
					X:        midX,
					Y:        midY,
					Interval: math.Abs(point2.X() - point1.X()),
					Distance: distance,
				})
			}
		}
	}

	return candidates
}

// getAdjustedPoint tries to move point farther from polygon edge
func getAdjustedPoint(x, y float64, path orb.Ring, vstep float64, weight func(orb.Point) float64) Candidate {
	p := Candidate{
		X: x,
		Y: y,
	}

	// Get distance to shape
	distance := getPointToShapeDistance(x, y, path)
	p.Distance = distance * weight(orb.Point{x, y})

	// Scan for better point (simplified)
	scanForBetterPoint(&p, path, vstep, weight)

	return p
}

// scanForBetterPoint tries to find a better point by scanning vertically
func scanForBetterPoint(p *Candidate, path orb.Ring, vstep float64, weight func(orb.Point) float64) {

	// Get the current candidate point
	currentPoint := orb.Point{p.X, p.Y}

	// Get the bounding box of the polygon
	bounds := path.Bound()

	// Start scanning from a bit below the current point
	minY := bounds.Min.Y() - vstep
	maxY := bounds.Max.Y() + vstep

	// Scan vertically from bottom to top
	for y := minY; y <= maxY; y += vstep {
		// Create a point at current y coordinate, using the same x as current point
		testPoint := orb.Point{currentPoint.X(), y}

		// Check if the point is inside the polygon
		if planar.RingContains(path, testPoint) {
			// Calculate the weight/distance for this point
			newWeight := weight(testPoint)

			// If this point has a better weight (lower distance, higher score, etc.)
			// Update the candidate if it's better than current
			if newWeight < p.Distance {
				p.X = testPoint.X()
				p.Y = testPoint.Y()
				p.Distance = newWeight
			}
		}
	}
}

// getPointToShapeDistance calculates distance from point to polygon
func getPointToShapeDistance(x, y float64, path orb.Ring) float64 {
	pt := orb.Point{x, y}
	return planar.DistanceFrom(path, pt)
}

// getInnerTics generates evenly spaced points between min and max
func getInnerTics(min, max float64, steps int) []float64 {
	if steps <= 0 {
		return []float64{}
	}

	rangeVal := max - min
	step := rangeVal / float64(steps+1)
	tics := make([]float64, 0, steps)

	for i := 1; i <= steps; i++ {
		tics = append(tics, min+step*float64(i))
	}

	return tics
}

// getBoundArea calculates area of a bound
func getBoundArea(bounds orb.Bound) float64 {
	width := getBoundWidth(bounds)
	height := getBoundHeight(bounds)
	return width * height
}

// getBoundWidth calculates width of a bound
func getBoundWidth(bounds orb.Bound) float64 {
	return bounds.Max[0] - bounds.Min[0]
}

// getBoundHeight calculates height of a bound
func getBoundHeight(bounds orb.Bound) float64 {
	return bounds.Max[1] - bounds.Min[1]
}

// Candidate represents a candidate point for anchor point calculation
type Candidate struct {
	X        float64
	Y        float64
	Interval float64
	Distance float64
}
