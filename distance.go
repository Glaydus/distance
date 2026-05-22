package distance

import "math"

// Coord represents a geographic coordinate with latitude and longitude.
type Coord struct {
	Lat float64
	Lon float64
}

// DistanceTo calculates the distance to another coordinate in meters.
func (c Coord) DistanceTo(other Coord) float64 {
	distance, _, _ := DistanceBetween(c.Lat, c.Lon, other.Lat, other.Lon)
	return distance
}

// ComputeDistanceBetween calculates the distance between two coordinates in meters.
func ComputeDistanceBetween(startCoord, endCoord Coord) float64 {
	distance, _, _ := DistanceBetween(startCoord.Lat, startCoord.Lon, endCoord.Lat, endCoord.Lon)
	return distance
}

// DistanceBetween calculates the approximate distance between two locations
// in meters using the WGS84 ellipsoid.
// The distance is calculated using the Vincenty formula, which is an iterative method
// using the Inverse Formula (section 4) from http://www.ngs.noaa.gov/PUBS_LIB/inverse.pdf.
// Based on the Android implementation of computeDistanceAndBearing.
func DistanceBetween(startLatitude, startLongitude, endLatitude, endLongitude float64) (distance float64, mi, km int) {

	const MAXITERS = 20
	// Convert lat/long to radians
	lat1 := startLatitude * math.Pi / 180
	lat2 := endLatitude * math.Pi / 180
	lon1 := startLongitude * math.Pi / 180
	lon2 := endLongitude * math.Pi / 180

	// WGS84 ellipsoid parameters
	a := 6378137.0      // WGS84 semi-major axis
	b := 6356752.314245 // WGS84 semi-minor axis
	f := (a - b) / a    // WGS84 flattening
	aSqMinusBSqOverBSq := (a*a - b*b) / (b * b)

	ll := lon2 - lon1
	u1 := math.Atan((1.0 - f) * math.Tan(lat1))
	u2 := math.Atan((1.0 - f) * math.Tan(lat2))

	cosU1 := math.Cos(u1)
	cosU2 := math.Cos(u2)
	sinU1 := math.Sin(u1)
	sinU2 := math.Sin(u2)
	cosU1cosU2 := cosU1 * cosU2
	sinU1sinU2 := sinU1 * sinU2

	var (
		A          float64
		sigma      float64
		deltaSigma float64
		cosSqAlpha float64
		cos2SM     float64
		cosSigma   float64
		sinSigma   float64
		cosLambda  float64
		sinLambda  float64
	)

	lambda := ll // initial approximation
	for range MAXITERS {
		lambdaOrig := lambda
		cosLambda = math.Cos(lambda)
		sinLambda = math.Sin(lambda)
		t1 := cosU2 * sinLambda
		t2 := cosU1*sinU2 - sinU1*cosU2*cosLambda
		sinSqSigma := t1*t1 + t2*t2 // (14)
		sinSigma = math.Sqrt(sinSqSigma)
		cosSigma = sinU1sinU2 + cosU1cosU2*cosLambda                       // (15)
		sigma = math.Atan2(sinSigma, cosSigma)                             // (16)
		sinAlpha := iif(sinSigma == 0, 0.0, cosU1cosU2*sinLambda/sinSigma) // (17)
		cosSqAlpha = 1.0 - sinAlpha*sinAlpha
		cos2SM = iif(cosSqAlpha == 0, float64(0.0), cosSigma-2.0*sinU1sinU2/cosSqAlpha) // (18)

		uSquared := cosSqAlpha * aSqMinusBSqOverBSq // definition
		A = 1 + (uSquared/16384.0)*                 // (3)
			(4096.0+uSquared*(-768+uSquared*(320.0-175.0*uSquared)))
		B := (uSquared / 1024.0) * // (4)
			(256.0 + uSquared*(-128.0+uSquared*(74.0-47.0*uSquared)))
		C := (f / 16.0) *
			cosSqAlpha * (4.0 + f*(4.0-3.0*cosSqAlpha)) // (10)
		cos2SMSq := cos2SM * cos2SM
		deltaSigma = B * sinSigma * // (6)
			(cos2SM + (B/4.0)*
				(cosSigma*(-1.0+2.0*cos2SMSq)-
					(B/6.0)*cos2SM*
						(-3.0+4.0*sinSigma*sinSigma)*
						(-3.0+4.0*cos2SMSq)))
		lambda = ll +
			(1.0-C)*f*sinAlpha*
				(sigma+C*sinSigma*
					(cos2SM+C*cosSigma*
						(-1.0+2.0*cos2SM*cos2SM))) // (11)
		delta := (lambda - lambdaOrig) / lambda
		if math.Abs(delta) < 1.0e-12 {
			break
		}
	}
	distance = (b * A * (sigma - deltaSigma))
	mi = roundInt(distance * 0.000621371192)
	km = roundInt(distance / 1000)
	return
}

func iif(cond bool, a, b float64) float64 {
	if cond {
		return a
	}
	return b
}

func roundInt(f float64) int {
	if f < 0 {
		return int(f - .5)
	}
	return int(f + .5)
}
