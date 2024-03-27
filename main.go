package main

import (
	"context"
	"strconv"
	"time"

	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/broadcast"
	"github.com/odin-software/metro/internal/models"

	"github.com/VividCortex/multitick"
	City "github.com/odin-software/metro/websites/city"
	Reporter "github.com/odin-software/metro/websites/reporter"
	VirtualWorld "github.com/odin-software/metro/websites/virtual-world"
)

var stationHashFunction = func(station models.Station) string {
	return strconv.FormatInt(station.ID, 10)
}

func main() {
	// Setup
	tick := multitick.NewTicker(20*time.Millisecond, -1*time.Millisecond)
	// tickerMap := time.NewTicker(1000 * time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	arrivals := make(chan broadcast.ADMessage[models.Train])
	departures := make(chan broadcast.ADMessage[models.Train])
	bcArr := broadcast.NewBroadcastServer(ctx, arrivals)
	bcDep := broadcast.NewBroadcastServer(ctx, departures)

	// Filling graph data.
	g := models.NewNetwork(stationHashFunction)
	sts, lines := GenerateTestData(bcArr, bcDep)
	// g.InsertVertices(sts)
	g.InsertVertices2(sts)
	g.InsertEdge(*sts[0], *sts[1], []models.Vector{models.NewVector(50.0, 250.0), models.NewVector(150.0, 200.0)})
	g.InsertEdge(*sts[1], *sts[2], []models.Vector{models.NewVector(250.0, 100.0)})
	g.InsertEdge(*sts[1], *sts[5], []models.Vector{models.NewVector(300.0, 300.0)})
	g.InsertEdge(*sts[1], *sts[3], []models.Vector{models.NewVector(350.0, 200.0), models.NewVector(400.0, 150.0), models.NewVector(400.0, 50.0)})
	g.InsertEdge(*sts[3], *sts[4], []models.Vector{models.NewVector(550.0, 100.0), models.NewVector(600.0, 100.0)})
	g.InsertEdge(*sts[3], *sts[10], []models.Vector{})
	g.InsertEdge(*sts[3], *sts[11], []models.Vector{models.NewVector(600.0, 50.0)})
	g.InsertEdge(*sts[5], *sts[6], []models.Vector{models.NewVector(100.0, 500.0)})
	g.InsertEdge(*sts[7], *sts[8], []models.Vector{models.NewVector(500.0, 450.0)})
	g.InsertEdge(*sts[8], *sts[9], []models.Vector{models.NewVector(500.0, 250.0), models.NewVector(550.0, 200.0)})

	// Creating the train and queing some destinations.
	chu4 := models.NewMake("4-Legged-chu", "A type of fast train.", 0.003, 1)
	chu1 := models.NewMake("1-Legged-chu", "Another type of fast train.", 0.004, 0.7)
	trains := make([]models.Train, 0)
	train2 := models.NewTrain("Cha", chu4, sts[0].Position, *sts[0], lines[0], &g, arrivals, departures)
	train3 := models.NewTrain("Che", chu1, sts[3].Position, *sts[3], lines[3], &g, arrivals, departures)
	train4 := models.NewTrain("Chi", chu1, sts[1].Position, *sts[11], lines[0], &g, arrivals, departures)
	train5 := models.NewTrain("Cho", chu4, sts[7].Position, *sts[7], lines[2], &g, arrivals, departures)
	train := models.NewTrain("Chu", chu1, sts[1].Position, *sts[1], lines[1], &g, arrivals, departures)
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
	baso := baso.NewBaso()
	bs := baso.ListStations()
	for _, st := range bs {
		println(st.Name)
	}

	go Reporter.ReporterServer()
	go VirtualWorld.VirtualWorldServer()
	City.CityServer()
}
