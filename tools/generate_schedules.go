package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	pixelsPerMeter      = 0.01    // 1 pixel = 100 meters
	ticksPerSecond      = 60.0    // 60 ticks per second
	dwellTimeSeconds    = 45      // Time trains spend at each station
	startTimeSeconds    = 28800   // 8:00 AM in seconds since midnight
)

type Station struct {
	ID   int64
	Name string
	X    float64
	Y    float64
}

type Train struct {
	ID           int64
	Name         string
	LineID       int64
	TopSpeed     float64 // pixels/tick
	Acceleration float64 // pixels/tick²
}

type ScheduleEntry struct {
	TrainID        int64
	StationID      int64
	ScheduledTime  int // seconds since midnight
	SequenceOrder  int
}

func main() {
	db, err := sql.Open("sqlite3", "metro.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Get all trains
	trains, err := getTrains(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting trains: %v\n", err)
		os.Exit(1)
	}

	// Get all stations
	stations, err := getStations(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting stations: %v\n", err)
		os.Exit(1)
	}

	// Generate schedules for each train
	var schedules []ScheduleEntry

	for _, train := range trains {
		// Get the line's stations in order
		lineStations, err := getLineStations(db, train.LineID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting line stations for train %s: %v\n", train.Name, err)
			continue
		}

		// Generate schedule for this train
		trainSchedules := generateTrainSchedule(train, lineStations, stations)
		schedules = append(schedules, trainSchedules...)
	}

	// Output SQL INSERT statements
	fmt.Println("-- Generated schedule data")
	fmt.Println("-- Assumes simulation starts at 8:00 AM (28800 seconds since midnight)")
	fmt.Println("-- Dwell time at each station: 45 seconds")
	fmt.Println()

	for _, sched := range schedules {
		hours := sched.ScheduledTime / 3600
		minutes := (sched.ScheduledTime % 3600) / 60
		seconds := sched.ScheduledTime % 60

		fmt.Printf("INSERT INTO schedule (train_id, station_id, scheduled_time, sequence_order) VALUES (%d, %d, %d, %d); -- %02d:%02d:%02d\n",
			sched.TrainID, sched.StationID, sched.ScheduledTime, sched.SequenceOrder, hours, minutes, seconds)
	}
}

func getTrains(db *sql.DB) ([]Train, error) {
	rows, err := db.Query(`
		SELECT t.id, t.name, t.lineId, m.top_speed, m.acceleration
		FROM train t
		JOIN make m ON t.makeId = m.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trains []Train
	for rows.Next() {
		var t Train
		if err := rows.Scan(&t.ID, &t.Name, &t.LineID, &t.TopSpeed, &t.Acceleration); err != nil {
			return nil, err
		}
		trains = append(trains, t)
	}
	return trains, rows.Err()
}

func getStations(db *sql.DB) (map[int64]Station, error) {
	rows, err := db.Query("SELECT id, name, x, y FROM station")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stations := make(map[int64]Station)
	for rows.Next() {
		var s Station
		if err := rows.Scan(&s.ID, &s.Name, &s.X, &s.Y); err != nil {
			return nil, err
		}
		stations[s.ID] = s
	}
	return stations, rows.Err()
}

func getLineStations(db *sql.DB, lineID int64) ([]int64, error) {
	rows, err := db.Query(`
		SELECT stationId FROM station_line
		WHERE lineId = ?
		ORDER BY odr ASC
	`, lineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stationIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		stationIDs = append(stationIDs, id)
	}
	return stationIDs, rows.Err()
}

func generateTrainSchedule(train Train, lineStations []int64, stations map[int64]Station) []ScheduleEntry {
	var schedules []ScheduleEntry
	currentTime := startTimeSeconds

	// Loop through the line multiple times to simulate a full day
	// Trains run from 8:00 AM to 10:00 PM (14 hours = 50400 seconds)
	endTime := startTimeSeconds + (14 * 3600) // 10:00 PM

	for currentTime < endTime {
		// Make one complete loop of the line
		for i := 0; i < len(lineStations); i++ {
			stationID := lineStations[i]

			schedules = append(schedules, ScheduleEntry{
				TrainID:       train.ID,
				StationID:     stationID,
				ScheduledTime: currentTime,
				SequenceOrder: len(schedules) + 1,
			})

			// Add dwell time
			currentTime += dwellTimeSeconds

			// Calculate travel time to next station (if not last station)
			if i < len(lineStations)-1 {
				nextStationID := lineStations[i+1]
				travelTime := calculateTravelTime(
					stations[stationID],
					stations[nextStationID],
					train.TopSpeed,
					train.Acceleration,
				)
				currentTime += travelTime
			}
		}

		// After completing the line, return to the first station
		// (Calculate travel time from last to first station)
		if len(lineStations) > 0 {
			lastStationID := lineStations[len(lineStations)-1]
			firstStationID := lineStations[0]
			travelTime := calculateTravelTime(
				stations[lastStationID],
				stations[firstStationID],
				train.TopSpeed,
				train.Acceleration,
			)
			currentTime += travelTime
		}
	}

	return schedules
}

func calculateTravelTime(from, to Station, topSpeed, acceleration float64) int {
	// Calculate distance in pixels
	dx := to.X - from.X
	dy := to.Y - from.Y
	distancePixels := math.Sqrt(dx*dx + dy*dy)

	// Convert to meters
	distanceMeters := distancePixels / pixelsPerMeter

	// Estimate travel time assuming:
	// 1. Acceleration phase to top speed (or half distance, whichever comes first)
	// 2. Constant speed phase
	// 3. Deceleration phase (same as acceleration)

	// Convert top speed from pixels/tick to m/s
	pixelsPerSecond := topSpeed * ticksPerSecond
	metersPerSecond := pixelsPerSecond / pixelsPerMeter

	// Convert acceleration from pixels/tick² to m/s²
	accelerationMPS2 := (acceleration * ticksPerSecond * ticksPerSecond) / pixelsPerMeter

	// Time and distance to reach top speed
	timeToTopSpeed := metersPerSecond / accelerationMPS2
	distanceToTopSpeed := 0.5 * accelerationMPS2 * timeToTopSpeed * timeToTopSpeed

	var totalTime float64

	if distanceToTopSpeed*2 >= distanceMeters {
		// Short distance: accelerate to midpoint, then decelerate
		// Total distance = 2 * (0.5 * a * t²)
		// Solve for t: distanceMeters = a * t²
		t := math.Sqrt(distanceMeters / accelerationMPS2)
		totalTime = t * 2 // Accel + decel
	} else {
		// Long distance: accelerate, cruise, decelerate
		cruiseDistance := distanceMeters - (distanceToTopSpeed * 2)
		cruiseTime := cruiseDistance / metersPerSecond
		totalTime = timeToTopSpeed + cruiseTime + timeToTopSpeed
	}

	return int(math.Ceil(totalTime))
}
