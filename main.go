package main

import (
	"context"
	"internal/baso"
	"internal/model"
	"time"

	"internal/broadcast"

	"github.com/VividCortex/multitick"
	City "github.com/odin-software/metro/websites/city"
	Reporter "github.com/odin-software/metro/websites/reporter"
	TwoD "github.com/odin-software/metro/websites/two-d"
	VirtualWorld "github.com/odin-software/metro/websites/virtual-world"
)

var stationHashFunction = func(station model.Station) string {
	return station.ID
}

func main() {
	// Setup
	tick := multitick.NewTicker(20*time.Millisecond, -1*time.Millisecond)
	// tickerMap := time.NewTicker(1000 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	arrivals := make(chan broadcast.ADMessage[model.Train])
	departures := make(chan broadcast.ADMessage[model.Train])
	bcArr := broadcast.NewBroadcastServer(ctx, arrivals)
	bcDep := broadcast.NewBroadcastServer(ctx, departures)

	// Filling graph data.
	g := model.NewNetwork[model.Station](stationHashFunction)
	sts, lines := GenerateTestData(bcArr, bcDep)
	// g.InsertVertices(sts)
	g.InsertVertices2(sts)
	g.InsertEdge(*sts[0], *sts[1], []model.Vector{model.NewVector(50.0, 250.0), model.NewVector(150.0, 200.0)})
	g.InsertEdge(*sts[1], *sts[2], []model.Vector{model.NewVector(250.0, 100.0)})
	g.InsertEdge(*sts[1], *sts[5], []model.Vector{model.NewVector(300.0, 300.0)})
	g.InsertEdge(*sts[1], *sts[3], []model.Vector{model.NewVector(350.0, 200.0), model.NewVector(400.0, 150.0), model.NewVector(400.0, 50.0)})
	g.InsertEdge(*sts[3], *sts[4], []model.Vector{model.NewVector(550.0, 100.0), model.NewVector(600.0, 100.0)})
	g.InsertEdge(*sts[3], *sts[10], []model.Vector{})
	g.InsertEdge(*sts[3], *sts[11], []model.Vector{model.NewVector(600.0, 50.0)})
	g.InsertEdge(*sts[5], *sts[6], []model.Vector{model.NewVector(100.0, 500.0)})
	g.InsertEdge(*sts[7], *sts[8], []model.Vector{model.NewVector(500.0, 450.0)})
	g.InsertEdge(*sts[8], *sts[9], []model.Vector{model.NewVector(500.0, 250.0), model.NewVector(550.0, 200.0)})

	// Creating the train and queing some destinations.
	chu4 := model.NewMake("4-Legged-chu", "A type of fast train.", 0.003, 1)
	chu1 := model.NewMake("1-Legged-chu", "Another type of fast train.", 0.004, 0.7)
	trains := make([]model.Train, 0)
	train := model.NewTrain("Chu", chu1, sts[1].Position, *sts[1], lines[1], &g, arrivals, departures)
	train2 := model.NewTrain("Cha", chu4, sts[0].Position, *sts[0], lines[0], &g, arrivals, departures)
	train3 := model.NewTrain("Che", chu1, sts[3].Position, *sts[3], lines[3], &g, arrivals, departures)
	train4 := model.NewTrain("Chi", chu1, sts[1].Position, *sts[11], lines[0], &g, arrivals, departures)
	train5 := model.NewTrain("Cho", chu4, sts[7].Position, *sts[7], lines[2], &g, arrivals, departures)
	trains = append(trains, train, train2, train3, train4, train5)

	// Starting the goroutines for the trains.
	// This should be changed eventually to have just one tick and then on tick call all the updates on goroutines.

	for i := 0; i < len(trains); i++ {
		go func(idx int) {
			sub := tick.Subscribe()
			for range sub {
				trains[idx].Update()
			}
		}(i)
	}

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
	baso.Stations()

	go Reporter.ReporterServer()
	go TwoD.TwoDimensionalServer()
	go VirtualWorld.VirtualWorldServer()
	City.CityServer()
}
