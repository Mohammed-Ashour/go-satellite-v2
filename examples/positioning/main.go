package main

import (
	"fmt"
	"time"

	"github.com/Mohammed-Ashour/go-satellite-v2/pkg/satellite"
	"github.com/Mohammed-Ashour/go-satellite-v2/pkg/tle"
)

func main() {
	tles, err := tle.ReadTLEFile("examples/tle/tle_example.txt")
	if err != nil {
		panic(err)
	}
	sat := satellite.NewSatelliteFromTLE(tles[0], satellite.GravityWGS84)
	fmt.Println("Satellite Name:", tles[0].Name)
	fmt.Println("Satellite Line 1:", tles[0].Line1)
	fmt.Println("Satellite Line 2:", tles[0].Line2)
	tle_time, err := tles[0].Time()
	if err != nil {
		panic(err)
	}
	tle_lat, tle_lon, tle_alt, tle_vel := sat.Locate(tle_time)
	fmt.Println("Satellite Epoch:", tle_time)

	fmt.Printf("TLE Latitude: %f, Longitude: %f, Altitude: %f, Velocity: %f\n", tle_lat, tle_lon, tle_alt, tle_vel)
	fmt.Println("Current UTC Time:", time.Now().UTC())
	t := time.Now().UTC()
	lat, lon, alt, vel := sat.Locate(t)
	fmt.Printf("Latitude: %f, Longitude: %f, Altitude: %f, Velocity: %f\n", lat, lon, alt, vel)
}
