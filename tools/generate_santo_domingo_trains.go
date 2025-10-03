package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("=== Santo Domingo Train Generator ===\n")

	db, err := sql.Open("sqlite3", "../data/metro.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Get line IDs
	linea1ID, linea2ID, err := getLineIDs(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting line IDs: %v\n", err)
		os.Exit(1)
	}

	// Get stations for each line
	linea1Stations, err := getLineStations(db, linea1ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting Línea 1 stations: %v\n", err)
		os.Exit(1)
	}

	linea2Stations, err := getLineStations(db, linea2ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting Línea 2 stations: %v\n", err)
		os.Exit(1)
	}

	// Get make ID (train type)
	makeID, err := getMakeID(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting make ID: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d stations on Línea 1\n", len(linea1Stations))
	fmt.Printf("Found %d stations on Línea 2\n", len(linea2Stations))
	fmt.Printf("Using make ID: %d\n\n", makeID)

	// Generate SQL seed file
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("../data/sql/seeds/%s_santo_domingo_trains.sql", timestamp)

	f, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	w := func(s string) {
		fmt.Fprint(f, s)
	}

	w("-- +goose Up\n")
	w("-- +goose StatementBegin\n")
	w("SELECT 'UP Santo Domingo trains seed';\n")
	w(fmt.Sprintf("-- Generated: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	w("-- Real-life fleet: Línea 1 = 40 trains, Línea 2 = 29 trains\n\n")

	trainID := 2000 // Start from high ID to avoid conflicts

	// Generate 40 trains for Línea 1
	// Random distribution of forward/backward directions
	w("-- Línea 1 Trains (40 trains - random directions)\n")
	for i := 0; i < 40; i++ {
		stationIdx := i % len(linea1Stations)
		station := linea1Stations[stationIdx]

		// Randomly assign forward (1) or backward (0) direction
		forward := i % 2 // Simple alternation for even distribution

		trainName := fmt.Sprintf("L1-T%02d", i+1)
		w(fmt.Sprintf("INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId, forward) VALUES (%d, '%s', %.2f, %.2f, 0.0, %d, %d, %d, %d);\n",
			trainID, trainName, station.X, station.Y, station.ID, makeID, linea1ID, forward))
		trainID++
	}
	w("\n")

	// Generate 29 trains for Línea 2
	// Random distribution of forward/backward directions
	w("-- Línea 2 Trains (29 trains - random directions)\n")
	for i := 0; i < 29; i++ {
		stationIdx := i % len(linea2Stations)
		station := linea2Stations[stationIdx]

		// Randomly assign forward (1) or backward (0) direction
		forward := i % 2 // Simple alternation for even distribution

		trainName := fmt.Sprintf("L2-T%02d", i+1)
		w(fmt.Sprintf("INSERT OR IGNORE INTO train (id, name, x, y, z, currentId, makeId, lineId, forward) VALUES (%d, '%s', %.2f, %.2f, 0.0, %d, %d, %d, %d);\n",
			trainID, trainName, station.X, station.Y, station.ID, makeID, linea2ID, forward))
		trainID++
	}
	w("\n")

	w("-- +goose StatementEnd\n\n")
	w("-- +goose Down\n")
	w("-- +goose StatementBegin\n")
	w("SELECT 'DOWN Santo Domingo trains seed';\n")
	w("DELETE FROM train WHERE id >= 2000 AND id < 3000;\n")
	w("-- +goose StatementEnd\n")

	fmt.Printf("✓ Generated: %s\n", filename)
	fmt.Println("\nCreated:")
	fmt.Println("  - 40 trains for Línea 1")
	fmt.Println("  - 29 trains for Línea 2")
	fmt.Println("  - Total: 69 trains")
	fmt.Println("\nTo apply:")
	fmt.Println("  1. make seed_santo_domingo  (loads trains)")
	fmt.Println("  2. cd tools && go run generate_schedules.go > schedules.sql")
	fmt.Println("  3. sqlite3 ../data/metro.db < schedules.sql")
}

type Station struct {
	ID   int64
	Name string
	X    float64
	Y    float64
}

func getLineIDs(db *sql.DB) (int64, int64, error) {
	var linea1ID, linea2ID int64

	err := db.QueryRow("SELECT id FROM line WHERE name = 'Línea 1'").Scan(&linea1ID)
	if err != nil {
		return 0, 0, fmt.Errorf("Línea 1 not found: %v", err)
	}

	err = db.QueryRow("SELECT id FROM line WHERE name = 'Línea 2'").Scan(&linea2ID)
	if err != nil {
		return 0, 0, fmt.Errorf("Línea 2 not found: %v", err)
	}

	return linea1ID, linea2ID, nil
}

func getLineStations(db *sql.DB, lineID int64) ([]Station, error) {
	rows, err := db.Query(`
		SELECT s.id, s.name, s.x, s.y
		FROM station s
		JOIN station_line sl ON s.id = sl.stationId
		WHERE sl.lineId = ?
		ORDER BY sl.odr ASC
	`, lineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stations []Station
	for rows.Next() {
		var s Station
		if err := rows.Scan(&s.ID, &s.Name, &s.X, &s.Y); err != nil {
			return nil, err
		}
		stations = append(stations, s)
	}
	return stations, rows.Err()
}

func getMakeID(db *sql.DB) (int64, error) {
	var makeID int64
	// Get the first available make (train type)
	err := db.QueryRow("SELECT id FROM make LIMIT 1").Scan(&makeID)
	return makeID, err
}
