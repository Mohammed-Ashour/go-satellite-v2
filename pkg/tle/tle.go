package tle

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type TLELine1 struct {
	LineNumber       string
	SataliteID       string
	Classification   string
	LaunchYear       string
	LaunchNumber     string
	LaunchPiece      string
	EpochYear        string
	EpochDay         string
	FirstDerivative  string
	SecondDerivative string
	Bstar            string // drag term or radiation pressure term
	EphemerisType    string
	ElementSetNumber string
	Checksum         string
	LineString       string
}
type TLELine2 struct {
	LineNumber        string
	SataliteID        string
	Inclination       string // degrees
	RightAscension    string // degrees
	Eccentricity      string
	ArgumentOfPerigee string // degrees
	MeanAnomaly       string // degrees
	MeanMotion        string // revolutions per day
	RevolutionNumber  string
	Checksum          string
	LineString        string
}

/*
ISS (ZARYA)
1 25544U 98067A   08264.51782528 -.00002182  00000-0 -11606-4 0  2927
2 25544  51.6416 247.4627 0006703 130.5360 325.0288 15.72125391563537
*/
type TLE struct {
	Name    string
	NoradID string
	Line1   TLELine1
	Line2   TLELine2
}

func (t TLE) String() string {
	return fmt.Sprintf("%s\n%s\n%s", t.Name, t.Line1.LineString, t.Line2.LineString)
}
func ReadTLELine1(line string) (TLELine1, error) {
	if len(line) < 69 {
		return TLELine1{}, fmt.Errorf("line 1 too short: %d chars", len(line))
	}

	tleLine1 := TLELine1{
		LineString: line,
	}

	// Fixed-width field parsing based on TLE format specification
	fields := map[string][2]int{
		"LineNumber":       {0, 1},
		"SatelliteID":      {2, 7},
		"Classification":   {7, 8},
		"LaunchYear":       {9, 11},
		"LaunchNumber":     {11, 14},
		"LaunchPiece":      {14, 17},
		"EpochYear":        {18, 20},
		"EpochDay":         {20, 32},
		"FirstDerivative":  {33, 43},
		"SecondDerivative": {44, 52},
		"Bstar":            {53, 61},
		"EphemerisType":    {62, 63},
		"ElementSetNumber": {64, 68},
		"Checksum":         {68, 69},
	}

	var err error
	for field, pos := range fields {
		value := strings.TrimSpace(line[pos[0]:pos[1]])
		switch field {
		case "SatelliteID":
			tleLine1.SataliteID = value
		case "Classification":
			tleLine1.Classification = value
		case "LaunchYear":
			tleLine1.LaunchYear = value
		case "LaunchNumber":
			tleLine1.LaunchNumber = value
		case "LaunchPiece":
			tleLine1.LaunchPiece = value
		case "EpochYear":
			tleLine1.EpochYear = value
		case "EpochDay":
			tleLine1.EpochDay = value
		case "FirstDerivative":
			tleLine1.FirstDerivative = ParseScientificNotation(value)
		case "SecondDerivative":
			tleLine1.SecondDerivative = ParseScientificNotation(value)
		case "Bstar":
			tleLine1.Bstar = ParseScientificNotation(value)
		case "EphemerisType":
			tleLine1.EphemerisType = value
		case "ElementSetNumber":
			tleLine1.ElementSetNumber = value
		case "Checksum":
			tleLine1.Checksum = value
		}
	}

	return tleLine1, err
}

func ReadTLELine2(line string) (TLELine2, error) {
	if len(line) < 69 {
		return TLELine2{}, fmt.Errorf("line 2 too short: %d chars", len(line))
	}

	tleLine2 := TLELine2{
		LineString: line,
	}

	// Fixed-width field parsing based on TLE format specification
	fields := map[string][2]int{
		"LineNumber":        {0, 1},
		"SatelliteID":       {2, 7},
		"Inclination":       {8, 16},
		"RightAscension":    {17, 25},
		"Eccentricity":      {26, 33},
		"ArgumentOfPerigee": {34, 42},
		"MeanAnomaly":       {43, 51},
		"MeanMotion":        {52, 63},
		"RevolutionNumber":  {63, 68},
		"Checksum":          {68, 69},
	}

	var err error
	for field, pos := range fields {
		value := strings.TrimSpace(line[pos[0]:pos[1]])
		switch field {
		case "SatelliteID":
			tleLine2.SataliteID = value
		case "Inclination":
			tleLine2.Inclination = value
		case "RightAscension":
			tleLine2.RightAscension = value
		case "Eccentricity":
			tleLine2.Eccentricity = "0." + value // Add leading "0." for eccentricity
		case "ArgumentOfPerigee":
			tleLine2.ArgumentOfPerigee = value
		case "MeanAnomaly":
			tleLine2.MeanAnomaly = value
		case "MeanMotion":
			tleLine2.MeanMotion = value
		case "RevolutionNumber":
			tleLine2.RevolutionNumber = value
		case "Checksum":
			tleLine2.Checksum = value
		}
	}

	return tleLine2, err
}
func ParseTLE(line1, line2, name string) (TLE, error) {
	tle := TLE{
		Name: name,
	}

	var err error
	tle.Line1, err = ReadTLELine1(line1)
	if err != nil {
		return TLE{}, err
	}

	tle.Line2, err = ReadTLELine2(line2)
	if err != nil {
		return TLE{}, err
	}

	return tle, nil
}

func ReadTLEFile(filePath string) ([]TLE, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var tles []TLE
	var currentTLE TLE

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "1 ") {
			currentTLE.Line1, err = ReadTLELine1(line)
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(line, "2 ") {
			components := strings.Fields(line)
			currentTLE.NoradID = components[1]
			currentTLE.Line2, err = ReadTLELine2(line)
			if err != nil {
				return nil, err
			}

			tles = append(tles, currentTLE)
			currentTLE = TLE{}
		} else {
			currentTLE.Name = strings.TrimSpace(line)
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tles, nil
}

// Time returns the time of the TLE epoch as a time.Time object.
func (t TLE) Time() (time.Time, error) {
	tleEpoch := fmt.Sprintf("%s%s", t.Line1.EpochYear, t.Line1.EpochDay)
	parts := strings.SplitN(tleEpoch, ".", 2)
	if len(parts) != 2 {
		return time.Time{}, fmt.Errorf("invalid TLE epoch format: '%s'. Expected format YYDDD.FFFFFFFF", tleEpoch)
	}
	yearDayPart := parts[0]
	fractionalDayPart := "0." + parts[1] // Prepend "0." for float parsing

	// Validate the length of the year/day part (must be 5 digits: YYDDD)
	if len(yearDayPart) != 5 {
		return time.Time{}, fmt.Errorf("invalid TLE epoch format: year/day part '%s' must be 5 digits (YYDDD)", yearDayPart)
	}

	// Extract YY and DDD
	yearStr := yearDayPart[:2]
	dayStr := yearDayPart[2:]

	// Parse the two-digit year (YY)
	yearYY, err := strconv.Atoi(yearStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid TLE epoch: cannot parse year '%s': %w", yearStr, err)
	}

	// Determine the full year based on the TLE convention
	var fullYear int
	if yearYY >= 57 { // Years 57-99 are 1957-1999
		fullYear = 1900 + yearYY
	} else { // Years 00-56 are 2000-2056
		fullYear = 2000 + yearYY
	}

	// Parse the day of the year (DDD)
	dayOfYear, err := strconv.Atoi(dayStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid TLE epoch: cannot parse day of year '%s': %w", dayStr, err)
	}

	// Basic validation for day of year
	if dayOfYear < 1 || dayOfYear > 366 { // Allow 366 for leap years
		return time.Time{}, fmt.Errorf("invalid TLE epoch: day of year %d out of range (1-366)", dayOfYear)
	}

	// Parse the fractional part of the day
	fractionalDay, err := strconv.ParseFloat(fractionalDayPart, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid TLE epoch: cannot parse fractional day '%s': %w", fractionalDayPart, err)
	}

	// Check fractional day bounds
	if fractionalDay < 0.0 || fractionalDay >= 1.0 {
		return time.Time{}, fmt.Errorf("invalid TLE epoch: fractional day %f out of range [0.0, 1.0)", fractionalDay)
	}

	// Calculate the start of the given year in UTC [[9]]
	startOfYear := time.Date(fullYear, 1, 1, 0, 0, 0, 0, time.UTC)

	// Add (dayOfYear - 1) days to get to the start of the specific day
	targetDayStart := startOfYear.AddDate(0, 0, dayOfYear-1)

	// Calculate the duration represented by the fractional day
	// Nanoseconds in a day = 24 hours * 60 minutes/hour * 60 seconds/minute * 1e9 nanoseconds/second
	nanosInDay := 24.0 * 60.0 * 60.0 * 1e9
	durationNanos := time.Duration(math.Round(fractionalDay * nanosInDay)) // Round to nearest nanosecond

	// Add the duration to the start of the target day
	finalTime := targetDayStart.Add(durationNanos)

	// Return the final time.Time object in UTC
	return finalTime.UTC(), nil
}
