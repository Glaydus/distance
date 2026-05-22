# distance

[![Go Report Card](https://goreportcard.com/badge/github.com/glaydus/distance)](https://goreportcard.com/report/github.com/glaydus/distance)
[![Release](https://img.shields.io/github/v/release/glaydus/distance)](https://github.com/glaydus/distance/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/glaydus/distance.svg)](https://pkg.go.dev/github.com/glaydus/distance)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/glaydus/distance)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A Go package for calculating geodesic distances between geographic coordinates using the **Vincenty Inverse Formula** on the WGS84 ellipsoid.

## Installation

```bash
go get github.com/glaydus/distance
```

## Usage

```go
import "github.com/glaydus/distance"
```

### DistanceBetween

The primary function. Takes two pairs of latitude/longitude coordinates and returns the distance in meters, miles, and kilometers.

```go
dist, mi, km := distance.DistanceBetween(52.2297, 21.0122, 50.0647, 19.9450)
fmt.Printf("Warsaw → Kraków: %.0f m  |  %d km  |  %d mi\n", dist, km, mi)
// Warsaw → Kraków: 252162 m  |  252 km  |  156 mi
```

### ComputeDistanceBetween

A convenience wrapper that accepts `Coord` structs and returns the distance in meters.

```go
warsaw := distance.Coord{Lat: 52.2297, Lon: 21.0122}
krakow := distance.Coord{Lat: 50.0647, Lon: 19.9450}

dist := distance.ComputeDistanceBetween(warsaw, krakow)
fmt.Printf("%.0f m\n", dist)
```

### Coord.DistanceTo

A method on `Coord` that returns the distance to another coordinate in meters.

```go
warsaw := distance.Coord{Lat: 52.2297, Lon: 21.0122}
krakow := distance.Coord{Lat: 50.0647, Lon: 19.9450}

dist := warsaw.DistanceTo(krakow)
fmt.Printf("%.0f m\n", dist)
```

## API Reference

### Types

```go
type Coord struct {
    Lat float64 // Latitude in decimal degrees
    Lon float64 // Longitude in decimal degrees
}
```

### Functions

| Signature | Returns | Description |
|-----------|---------|-------------|
| `DistanceBetween(startLat, startLon, endLat, endLon float64)` | `(distance float64, mi int, km int)` | Distance between two lat/lon pairs |
| `ComputeDistanceBetween(start, end Coord)` | `float64` | Distance between two `Coord` values in meters |
| `(c Coord) DistanceTo(other Coord)` | `float64` | Distance from receiver to another `Coord` in meters |

## Algorithm

### Why not the Haversine formula?

The commonly used **Haversine formula** treats the Earth as a perfect sphere. This introduces errors of up to ~0.5% for long distances, which can translate to several kilometers on intercontinental routes.

### Vincenty Inverse Formula

This package uses the **Vincenty Inverse Formula** (T. Vincenty, 1975), which models the Earth as an oblate spheroid — specifically the **WGS84 ellipsoid**, the same reference system used by GPS.

WGS84 parameters used:

| Parameter | Value | Description |
|-----------|-------|-------------|
| `a` | 6 378 137.0 m | Semi-major axis (equatorial radius) |
| `b` | 6 356 752.314245 m | Semi-minor axis (polar radius) |
| `f` | 1 / 298.257223563 | Flattening `(a − b) / a` |

The algorithm iterates up to 20 times, converging when the relative change in the longitude difference `λ` drops below `1×10⁻¹²`. In practice, convergence is reached in 2–3 iterations for most point pairs.

The implementation follows **Section 4 (Inverse Formula)** of:

> T. Vincenty, *"Direct and Inverse Solutions of Geodesics on the Ellipsoid with Application of Nested Equations"*, Survey Review, 1975.
> Reference implementation: [NGS Publication](http://www.ngs.noaa.gov/PUBS_LIB/inverse.pdf)

The code is based on the Android platform's `computeDistanceAndBearing` implementation, which uses the same formulation.

### Accuracy

Vincenty's formula is accurate to within **0.5 mm** on the WGS84 ellipsoid, making it suitable for virtually all practical applications — from navigation to geospatial analytics.

## Running Tests

```bash
go test ./...
```

With benchmarks:

```bash
go test -bench=. -benchmem ./...
```

## Requirements

- Go 1.22 or later (uses `for range N` loop syntax)

## License

MIT
