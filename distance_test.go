package distance

import (
	"math"
	"testing"
)

// tolerance for float64 comparisons (0.1% error)
const tolerance = 0.001

func withinTolerance(a, b, tol float64) bool {
	if b == 0 {
		return math.Abs(a) < tol
	}
	return math.Abs(a-b)/math.Abs(b) < tol
}

// --- DistanceBetween ---

func TestDistanceBetween_SamePoint(t *testing.T) {
	dist, mi, km := DistanceBetween(52.2297, 21.0122, 52.2297, 21.0122)
	if dist != 0 {
		t.Errorf("same point: expected 0, got %f", dist)
	}
	if mi != 0 {
		t.Errorf("same point: expected mi=0, got %d", mi)
	}
	if km != 0 {
		t.Errorf("same point: expected km=0, got %d", km)
	}
}

func TestDistanceBetween_KnownDistance(t *testing.T) {
	// Warsaw → Kraków, known distance ~252 km
	dist, mi, km := DistanceBetween(52.2297, 21.0122, 50.0647, 19.9450)

	expectedDist := 252162.0 // ~252 km in meters
	if !withinTolerance(dist, expectedDist, 0.02) {
		t.Errorf("Warsaw→Kraków: expected ~%f m, got %f m", expectedDist, dist)
	}
	if mi < 150 || mi > 165 {
		t.Errorf("Warsaw→Kraków: expected mi in range [150,165], got %d", mi)
	}
	if km < 248 || km > 256 {
		t.Errorf("Warsaw→Kraków: expected km in range [248,256], got %d", km)
	}
}

func TestDistanceBetween_Symmetry(t *testing.T) {
	// Distance A→B should equal B→A
	dist1, _, _ := DistanceBetween(52.2297, 21.0122, 50.0647, 19.9450)
	dist2, _, _ := DistanceBetween(50.0647, 19.9450, 52.2297, 21.0122)

	if !withinTolerance(dist1, dist2, tolerance) {
		t.Errorf("symmetry: dist(A→B)=%f != dist(B→A)=%f", dist1, dist2)
	}
}

func TestDistanceBetween_EquatorPoints(t *testing.T) {
	// Two points on the equator 1 degree of longitude apart ≈ 111 320 m
	dist, _, km := DistanceBetween(0, 0, 0, 1)

	expectedDist := 111320.0
	if !withinTolerance(dist, expectedDist, 0.01) {
		t.Errorf("equator 1°: expected ~%f m, got %f m", expectedDist, dist)
	}
	if km != 111 {
		t.Errorf("equator 1°: expected km=111, got %d", km)
	}
}

func TestDistanceBetween_LongDistance(t *testing.T) {
	// Warsaw → New York, ~6873 km
	dist, _, km := DistanceBetween(52.2297, 21.0122, 40.7128, -74.0060)

	if !withinTolerance(dist, 6873000, 0.05) {
		t.Errorf("Warsaw→NYC: expected ~6873000 m, got %f m", dist)
	}
	if km < 6500 || km > 7200 {
		t.Errorf("Warsaw→NYC: expected km in range [6500,7200], got %d", km)
	}
}

func TestDistanceBetween_NegativeCoordinates(t *testing.T) {
	// Buenos Aires → Sydney (both in the southern hemisphere)
	dist, _, km := DistanceBetween(-34.6037, -58.3816, -33.8688, 151.2093)

	if !withinTolerance(dist, 11800000, 0.05) {
		t.Errorf("Buenos Aires→Sydney: expected ~11800000 m, got %f m", dist)
	}
	if km < 11200 || km > 12400 {
		t.Errorf("Buenos Aires→Sydney: expected km in range [11200,12400], got %d", km)
	}
}

// --- ComputeDistanceBetween ---

func TestComputeDistanceBetween_SamePoint(t *testing.T) {
	c := Coord{Lat: 48.8566, Lon: 2.3522}
	dist := ComputeDistanceBetween(c, c)
	if dist != 0 {
		t.Errorf("same point: expected 0, got %f", dist)
	}
}

// --- Coord.DistanceTo ---

func TestCoordDistanceTo_Symmetry(t *testing.T) {
	a := Coord{Lat: 52.2297, Lon: 21.0122}
	b := Coord{Lat: 48.8566, Lon: 2.3522}

	if !withinTolerance(a.DistanceTo(b), b.DistanceTo(a), tolerance) {
		t.Errorf("DistanceTo symmetry: a→b=%f, b→a=%f", a.DistanceTo(b), b.DistanceTo(a))
	}
}

// --- roundInt ---

func TestRoundInt(t *testing.T) {
	cases := []struct {
		input    float64
		expected int
	}{
		{0.0, 0},
		{0.4, 0},
		{0.5, 1},
		{1.5, 2},
		{2.9, 3},
		{100.0, 100},
		{-0.4, 0},
		{-0.5, -1},
		{-1.5, -2},
		{-2.9, -3},
	}
	for _, tc := range cases {
		got := roundInt(tc.input)
		if got != tc.expected {
			t.Errorf("roundInt(%f): expected %d, got %d", tc.input, tc.expected, got)
		}
	}
}

// --- Benchmarks ---

func BenchmarkDistanceBetween(b *testing.B) {
	for b.Loop() {
		DistanceBetween(52.2297, 21.0122, 50.0647, 19.9450)
	}
}

func BenchmarkComputeDistanceBetween(b *testing.B) {
	start := Coord{Lat: 52.2297, Lon: 21.0122}
	end := Coord{Lat: 50.0647, Lon: 19.9450}
	for b.Loop() {
		ComputeDistanceBetween(start, end)
	}
}

func BenchmarkCoordDistanceTo(b *testing.B) {
	a := Coord{Lat: 52.2297, Lon: 21.0122}
	c := Coord{Lat: 50.0647, Lon: 19.9450}
	for b.Loop() {
		a.DistanceTo(c)
	}
}
