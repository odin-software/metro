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
	"github.com/odin-software/metro/internal/broadcast"
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
	mapTick := time.NewTicker(control.DefaultConfig.TerminalMapDuration)
	ctx, cancel := context.WithCancel(context.Background())
	control.InitLogger()
	defer cancel()

	// Creating the broadcast channels for the trains.
	arrivals := make(chan broadcast.ADMessage[models.Train])
	departures := make(chan broadcast.ADMessage[models.Train])
	bcArr := broadcast.NewBroadcastServer(ctx, arrivals)
	bcDep := broadcast.NewBroadcastServer(ctx, departures)

	// Creating the city graph.
	cityNetwork := models.NewNetwork(StationHashFunction)

	// Loading stations, lines, edges from the database.
	stations := data.LoadStations(bcArr, bcDep)
	lines := data.LoadLines()
	err := cityNetwork.InsertVertices(stations)
	if err != nil {
		return
	}
	data.LoadEdges(&cityNetwork)

	// Creating the train with lines and channels to communicate.
	trains := data.LoadTrains(stations, lines, &cityNetwork, arrivals, departures)

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

	// Drawing a map in the console of the trains and stations.
	if control.DefaultConfig.TerminalMapEnabled {
		go StartMap(mapTick.C, stations, trains)
	}

	// Reflect what's on memory on the DB.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range reflexTick.C {
			data.DumpTrainsData(trains)
		}
	}()

	game := display.NewGame(trains, stations)
	ebiten.SetWindowSize(
		control.DefaultConfig.DisplayScreenWidth,
		control.DefaultConfig.DisplayScreenHeight,
	)
	ebiten.SetWindowTitle("Metro")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}
