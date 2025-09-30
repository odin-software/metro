package main

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/data"
	"github.com/odin-software/metro/display"
	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/models"
	"github.com/odin-software/metro/internal/sematick"
	"github.com/odin-software/metro/internal/tenjin"
)

var StationHashFunction = func(station models.Station) string {
	return strconv.FormatInt(station.ID, 10)
}

func main() {
	// Setup.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	loopTick := sematick.NewTicker(
		control.DefaultConfig.LoopDuration,
		control.DefaultConfig.LoopStartingState,
	)
	reflexTick := time.NewTicker(control.DefaultConfig.ReflexDuration)
	spawnTick := time.NewTicker(control.DefaultConfig.PassengerSpawnRate)
	defer spawnTick.Stop()
	control.InitLogger()

	// Initialize database (create and run migrations if needed).
	if err := data.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Creating the city graph.
	cityNetwork := models.NewNetwork(StationHashFunction)

	// Loading stations, lines, edges from the database.
	stations := data.LoadStations()
	lines := data.LoadLines(stations) // Pass stations so lines reference same pointers
	err := cityNetwork.InsertVertices(stations)
	if err != nil {
		return
	}
	data.LoadEdges(&cityNetwork)

	// Initialize Tenjin (the brain) if enabled
	var brain *tenjin.Tenjin
	var eventChannel chan<- interface{}
	if control.DefaultConfig.TenjinEnabled {
		// Count trains using baso
		db := baso.NewBaso()
		trainsData := db.ListTrainsFull()
		trainCount := len(trainsData)

		brain, err = tenjin.NewTenjin(trainCount)
		if err != nil {
			log.Fatal("Failed to initialize Tenjin:", err)
		}
		eventChannel = brain.GetEventChannel()
		control.Log("Tenjin initialized successfully")
	}

	// Creating the train with lines.
	trains := data.LoadTrains(stations, lines, &cityNetwork, eventChannel)

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

	// Start Tenjin if enabled
	if control.DefaultConfig.TenjinEnabled && brain != nil {
		brain.Start()
		control.Log("Tenjin brain started")
	}

	// Start passenger spawning
	data.SpawnPassengers(ctx, &wg, stations, lines, spawnTick, eventChannel)

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

	// Set monitor if configured (0 = primary, 1+ = other monitors)
	monitors := ebiten.AppendMonitors(nil)
	if control.DefaultConfig.DisplayMonitor >= 0 && control.DefaultConfig.DisplayMonitor < len(monitors) {
		ebiten.SetMonitor(monitors[control.DefaultConfig.DisplayMonitor])
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	// Cleanup: Stop Tenjin gracefully
	if control.DefaultConfig.TenjinEnabled && brain != nil {
		brain.Stop()
	}

	wg.Wait()
}
