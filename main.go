package main

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/data"
	"github.com/odin-software/metro/display"
	"github.com/odin-software/metro/internal/models"
	"github.com/odin-software/metro/internal/sematick"
)

var StationHashFunction = func(station models.Station) string {
	return strconv.FormatInt(station.ID, 10)
}

func main() {
	// Setup.
	var wg sync.WaitGroup
	loopTick := sematick.NewTicker(
		control.DefaultConfig.LoopDuration,
		control.DefaultConfig.LoopStartingState,
	)
	reflexTick := time.NewTicker(control.DefaultConfig.ReflexDuration)
	control.InitLogger()

	// Initialize database (create and run migrations if needed).
	if err := data.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Creating the city graph.
	cityNetwork := models.NewNetwork(StationHashFunction)

	// Loading stations, lines, edges from the database.
	stations := data.LoadStations()
	lines := data.LoadLines()
	err := cityNetwork.InsertVertices(stations)
	if err != nil {
		return
	}
	data.LoadEdges(&cityNetwork)

	// Creating the train with lines.
	trains := data.LoadTrains(stations, lines, &cityNetwork)

	// Starting the goroutines for the trains.
	for i := range len(trains) {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			sub := loopTick.Subscribe()
			for range sub {
				trains[idx].Tick()
			}
		}(i)
	}

	// Reflect what's on memory on the DB.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range reflexTick.C {
			data.DumpTrainsData(trains)
		}
	}()

	game := display.NewGame(trains, stations, lines)
	ebiten.SetWindowSize(
		control.DefaultConfig.DisplayScreenWidth*2,
		control.DefaultConfig.DisplayScreenHeight*2,
	)
	ebiten.SetWindowTitle("Metro")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
