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

	// Loading stations and lines from the database.
	stations := LoadStations(bcArr, bcDep)
	lines := LoadLines()

	cityNetwork.InsertVertices2(stations)
	cityNetwork.InsertEdge(*stations[0], *stations[1], []models.Vector{models.NewVector(50.0, 250.0), models.NewVector(150.0, 200.0)})
	cityNetwork.InsertEdge(*stations[1], *stations[2], []models.Vector{models.NewVector(250.0, 100.0)})
	cityNetwork.InsertEdge(*stations[1], *stations[5], []models.Vector{models.NewVector(300.0, 300.0)})
	cityNetwork.InsertEdge(*stations[1], *stations[3], []models.Vector{models.NewVector(350.0, 200.0), models.NewVector(400.0, 150.0), models.NewVector(400.0, 50.0)})
	cityNetwork.InsertEdge(*stations[3], *stations[4], []models.Vector{models.NewVector(550.0, 100.0), models.NewVector(600.0, 100.0)})
	cityNetwork.InsertEdge(*stations[3], *stations[10], []models.Vector{})
	cityNetwork.InsertEdge(*stations[3], *stations[11], []models.Vector{models.NewVector(600.0, 50.0)})
	cityNetwork.InsertEdge(*stations[5], *stations[6], []models.Vector{models.NewVector(100.0, 500.0)})
	cityNetwork.InsertEdge(*stations[7], *stations[8], []models.Vector{models.NewVector(500.0, 450.0)})
	cityNetwork.InsertEdge(*stations[8], *stations[9], []models.Vector{models.NewVector(500.0, 250.0), models.NewVector(550.0, 200.0)})

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

	// Starting the server for The New Metro Times, Virtual World and CityServer.
	go Reporter.ReporterServer()
	go VirtualWorld.VirtualWorldServer()
	City.CityServer()
}
