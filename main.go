package main

import (
	"context"
	"strconv"
	"time"

	"github.com/odin-software/metro/internal/broadcast"
	"github.com/odin-software/metro/internal/models"

	City "github.com/odin-software/metro/websites/city"
	Reporter "github.com/odin-software/metro/websites/reporter"
	VirtualWorld "github.com/odin-software/metro/websites/virtual-world"
)

var stationHashFunction = func(station models.Station) string {
	return strconv.FormatInt(station.ID, 10)
}

func main() {
	// Setup
	loopTick := time.NewTicker(20 * time.Millisecond)
	quit := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Creating the broadcast channels for the trains.
	arrivals := make(chan broadcast.ADMessage[models.Train])
	departures := make(chan broadcast.ADMessage[models.Train])
	bcArr := broadcast.NewBroadcastServer(ctx, arrivals)
	bcDep := broadcast.NewBroadcastServer(ctx, departures)

	// Creating the city graph.
	cityNetwork := models.NewNetwork(stationHashFunction)

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
	go func() {
		for {
			select {
			case <-loopTick.C:
				for i := 0; i < len(trains); i++ {
					go trains[i].Update()
				}
			case <-quit:
				loopTick.Stop()
				return
			}
		}
	}()

	// Drawing a map in the console of the trains and stations.
	// for range tickerMap.C {
	// 	fmt.Println(len(sts[0].Trains))
	// 	fmt.Println(len(sts[1].Trains))
	// 	fmt.Println(len(sts[2].Trains))
	// 	fmt.Println(len(sts[3].Trains))
	// 	fmt.Println(len(sts[4].Trains))
	// 	go PrintMap(800, 600, sts, trains)
	// }

	// Starting the server for The New Metro Times.
	go Reporter.ReporterServer()
	go VirtualWorld.VirtualWorldServer()
	City.CityServer()
}
