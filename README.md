# go-satellite-v2
```go
import "github.com/Mohammed-Ashour/go-satellite-v2/pkg/satellite"
```
A fork of [github.com/joshuaferrara/go-satellite](https://github.com/joshuaferrara/go-satellite) with improvements and additional features.

## Intro

This is a Go implementation of the SGP4 satellite propagation library, originally ported by Joshua Ferrara. The SGP4 model is used to track satellites and space debris based on two-line element (TLE) data.

## Changes from Original

This fork includes:
- Added TLE parsing and validation abstraction
- Improved code organization and package structure
- Additional test coverage
- Bug fixes and performance improvements
- Better error handling

## Usage

```go
import (
    "github.com/Mohammed-Ashour/go-satellite-v2/pkg/satellite"
    "github.com/Mohammed-Ashour/go-satellite-v2/pkg/tle"
)
```

### TLE Handling

```go
// Parse TLE
tle, err := tle.ParseTLE(line1, line2, name)

// Get epoch time
epochTime, err := tle.Time()

// Read TLE from file
tles, err := tle.ReadTLEFile("path/to/tle.txt")

```

### Satellite Operations

```go
// Create satellite from TLE
sat := satellite.NewSatelliteFromTLE(tle, gravity)

// Get satellite position at specific time
lat, lon, alt, vel := sat.Locate(time.Now())
```

### Constants

```go
const DEG2RAD float64 = math.Pi / 180.0
const RAD2DEG float64 = 180.0 / math.Pi
const TWOPI float64 = math.Pi * 2.0
const XPDOTP float64 = 1440.0 / (2.0 * math.Pi)
```

### Time Functions

```go
// Calculate Julian Date
jday := satellite.JDay(year, mon, day, hr, min, sec)

// Calculate Greenwich Mean Sidereal Time
gmst := satellite.GSTimeFromDate(year, mon, day, hr, min, sec)

// Calculate GMST from Julian date
theta := satellite.ThetaG_JD(jday)
```

### Coordinate Conversions

```go
// ECI to Latitude/Longitude/Altitude
altitude, velocity, latlong := satellite.ECIToLLA(eciCoords, gmst)

// ECI to ECEF
ecfCoords := satellite.ECIToECEF(eciCoords, gmst)

// Latitude/Longitude/Altitude to ECI
eciCoords := satellite.LLAToECI(obsCoords, alt, jday)
```

### Types

#### TLE
```go
func LatLongDeg(rad LatLong) (deg LatLong)
```
Convert LatLong in radians to LatLong in degrees

#### type LookAngles

```go
type LookAngles struct {
	Az, El, Rg float64
}
```

Holds an azimuth, elevation and range

#### func  ECIToLookAngles

```go
func ECIToLookAngles(eciSat Vector3, obsCoords LatLong, obsAlt, jday float64) (lookAngles LookAngles)
```
Calculate look angles for given satellite position and observer position obsAlt
in km Reference: http://celestrak.com/columns/v02n02/

#### type Satellite

```go
type Satellite struct {
    Line1 string
    Line2 string
}
```

#### Vector3
```go
type Vector3 struct {
    X, Y, Z float64
}
```

#### LatLong
```go
type LatLong struct {
    Latitude, Longitude float64
}
```

#### LookAngles
```go
type LookAngles struct {
    Az, El, Rg float64
}
```

## Error Handling

The library now includes proper error handling for TLE parsing and validation:
- Invalid TLE format
- Checksum verification
- Date/time parsing
- Field validation

## License

This project is licensed under the same terms as the original repository.

## Acknowledgments

- Original implementation by Joshua Ferrara ([github.com/joshuaferrara/go-satellite](https://github.com/joshuaferrara/go-satellite))
- Based on the SGP4 satellite propagation algorithms

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
