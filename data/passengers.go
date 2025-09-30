package data

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/odin-software/metro/internal/models"
)

// SpawnPassengers creates initial passengers and spawns new ones periodically
func SpawnPassengers(
	ctx context.Context,
	wg *sync.WaitGroup,
	stations []*models.Station,
	lines []models.Line,
	spawnTick *time.Ticker,
	eventChannel chan<- interface{},
) {
	// Build map of station -> reachable destinations (stations on same lines)
	stationDestinations := buildStationDestinationMap(stations, lines)

	// Initial spawn: create passengers at each station
	for _, station := range stations {
		spawnPassengersAtStation(station, stationDestinations, 3, eventChannel) // 3 per station initially
	}

	// Random spawning loop
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-spawnTick.C:
				// Random station gets 1-2 new passengers
				if len(stations) > 0 {
					station := stations[rand.Intn(len(stations))]
					count := rand.Intn(2) + 1
					spawnPassengersAtStation(station, stationDestinations, count, eventChannel)
				}
			}
		}
	}()
}

// buildStationDestinationMap creates a map of station ID -> reachable destinations
func buildStationDestinationMap(stations []*models.Station, lines []models.Line) map[int64][]*models.Station {
	stationDestinations := make(map[int64][]*models.Station)

	// For each station, find all reachable destinations (stations on same lines)
	for _, station := range stations {
		reachableSet := make(map[int64]*models.Station)

		// Find all lines that serve this station
		for _, line := range lines {
			stationOnLine := false
			for _, lineStation := range line.Stations {
				if lineStation.ID == station.ID {
					stationOnLine = true
					break
				}
			}

			// If station is on this line, add all other stations from this line
			if stationOnLine {
				for _, dest := range line.Stations {
					if dest.ID != station.ID {
						reachableSet[dest.ID] = dest
					}
				}
			}
		}

		// Convert set to slice
		reachable := make([]*models.Station, 0, len(reachableSet))
		for _, dest := range reachableSet {
			reachable = append(reachable, dest)
		}
		stationDestinations[station.ID] = reachable
	}

	return stationDestinations
}

func spawnPassengersAtStation(
	station *models.Station,
	stationDestinations map[int64][]*models.Station,
	count int,
	eventChannel chan<- interface{},
) {
	reachableStations := stationDestinations[station.ID]
	if len(reachableStations) == 0 {
		fmt.Printf("WARNING: Station %s (ID:%d) has no reachable destinations\n", station.Name, station.ID)
		return
	}

	for i := 0; i < count; i++ {
		// Pick random reachable destination
		dest := reachableStations[rand.Intn(len(reachableStations))]

		id := fmt.Sprintf("P-%d-%d", time.Now().Unix(), rand.Intn(10000))
		name := fmt.Sprintf("Passenger-%d", rand.Intn(1000))

		passenger := models.NewPassenger(id, name, station, dest, eventChannel)
		station.AddPassenger(passenger)
	}
}
