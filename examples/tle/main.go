package main

import (
	"fmt"

	"github.com/Mohammed-Ashour/go-satellite-v2/pkg/tle"
)

func main() {
	tles, err := tle.ReadTLEFile("examples/tle/tle_example.txt")
	if err != nil {

		panic(err)
	}
	for _, t := range tles {

		fmt.Println(t)
	}
	fmt.Println("TLEs read successfully.")
	fmt.Println("Number of TLEs:", len(tles))
	fmt.Println("First TLE Name:", tles[0].Name)
	fmt.Println("First TLE Line 1:", tles[0].Line1)
	fmt.Println("First TLE Line 2:", tles[0].Line2)
	fmt.Printf("First TLE Epoch:%s%s\n", tles[0].Line1.EpochYear, tles[0].Line1.EpochDay)
}
