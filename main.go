package main

import (
	"context"
	"strconv"
	"time"

	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/data"
	"github.com/odin-software/metro/internal/broadcast"
	"github.com/odin-software/metro/internal/models"
	"github.com/odin-software/metro/internal/sematick"

	"github.com/odin-software/metro/websites/city"
	"github.com/odin-software/metro/websites/events"
	"github.com/odin-software/metro/websites/reporter"
	"github.com/odin-software/metro/websites/virtual-world"
)

var StationHashFunction = func(station models.Station) string {
	return strconv.FormatInt(station.ID, 10)
}

func main() {
	// Setup.
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
	for i := 0; i < len(trains); i++ {
		go func(idx int) {
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
	go func() {
		for range reflexTick.C {
			data.DumpTrainsData(trains)
		}
	}()

	// Starting the server for The New Metro Times, Virtual World and CityServer.
	go reporter.Server()
	go virtual.Server()
	go events.Main()
	city.Main(loopTick)
}
