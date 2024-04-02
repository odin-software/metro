package main

import (
	"context"
	"time"

	"github.com/odin-software/metro/internal/broadcast"
	"github.com/odin-software/metro/internal/models"

	"github.com/VividCortex/multitick"
	City "github.com/odin-software/metro/websites/city"
	Reporter "github.com/odin-software/metro/websites/reporter"
	VirtualWorld "github.com/odin-software/metro/websites/virtual-world"
)

func main() {
	// Setup
	loopTick := multitick.NewTicker(DefaultConfig.LoopDuration, DefaultConfig.LoopDurationOffset)
	reflexTick := time.NewTicker(DefaultConfig.ReflexDuration)
	mapTick := time.NewTicker(DefaultConfig.TerminalMapDuration)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Creating the broadcast channels for the trains.
	arrivals := make(chan broadcast.ADMessage[models.Train])
	departures := make(chan broadcast.ADMessage[models.Train])
	bcArr := broadcast.NewBroadcastServer(ctx, arrivals)
	bcDep := broadcast.NewBroadcastServer(ctx, departures)

	// Creating the city graph.
	cityNetwork := models.NewNetwork(StationHashFunction)

	// Loading stations, lines, edges from the database.
	stations := LoadStations(bcArr, bcDep)
	lines := LoadLines()
	cityNetwork.InsertVertices(stations)
	LoadEdges(&cityNetwork)

	// Creating the train with lines and channels to communicate.
	trains := LoadTrains(stations, lines, &cityNetwork, arrivals, departures)

	// Starting the goroutines for the trains.
	for i := 0; i < len(trains); i++ {
		go func(idx int) {
			sub := loopTick.Subscribe()
			for range sub {
				trains[idx].Update()
			}
		}(i)
	}

	// Drawing a map in the console of the trains and stations.
	if DefaultConfig.TerminalMapEnabled {
		StartMap(mapTick.C, stations, trains)
	}

	// Reflect what's on memory on the DB.
	go func() {
		for range reflexTick.C {
			DumpTrainsData(trains)
		}
	}()

	// Starting the server for The New Metro Times, Virtual World and CityServer.
	go Reporter.ReporterServer()
	go VirtualWorld.VirtualWorldServer()
	City.CityServer()
}
