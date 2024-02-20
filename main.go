package main

import (
	"internal/model"
	"time"

	"github.com/VividCortex/multitick"
)

var stationHashFunction = func(station model.Station) string {
	return station.ID
}

func main() {
	tick := multitick.NewTicker(15*time.Millisecond, -1*time.Millisecond)
	tickerMap := time.NewTicker(1000 * time.Millisecond)

	// Filling graph data.
	g := model.NewNetwork[model.Station](stationHashFunction)
	sts, lines := GenerateTestData()
	g.InsertVertices(sts)
	g.InsertEdge(sts[0], sts[1], []model.Vector{model.NewVector(50.0, 250.0), model.NewVector(150.0, 200.0)})
	g.InsertEdge(sts[1], sts[2], []model.Vector{model.NewVector(250.0, 100.0)})
	g.InsertEdge(sts[1], sts[5], []model.Vector{model.NewVector(300.0, 300.0)})
	g.InsertEdge(sts[1], sts[3], []model.Vector{model.NewVector(350.0, 200.0), model.NewVector(400.0, 150.0), model.NewVector(400.0, 50.0)})
	g.InsertEdge(sts[3], sts[4], []model.Vector{model.NewVector(550.0, 100.0), model.NewVector(600.0, 100.0)})
	g.InsertEdge(sts[3], sts[10], []model.Vector{})
	g.InsertEdge(sts[3], sts[11], []model.Vector{model.NewVector(600.0, 50.0)})
	g.InsertEdge(sts[5], sts[6], []model.Vector{model.NewVector(100.0, 500.0)})
	g.InsertEdge(sts[7], sts[8], []model.Vector{model.NewVector(500.0, 450.0)})
	g.InsertEdge(sts[8], sts[9], []model.Vector{model.NewVector(500.0, 250.0), model.NewVector(550.0, 200.0)})

	// Creating the train and queing some destinations.
	chu4 := model.NewMake("4-Legged-chu", "A type of fast train.", 0.01, 4)
	chu1 := model.NewMake("1-Legged-chu", "Another type of fast train.", 0.02, 7)
	trains := make([]model.Train, 0)
	train := model.NewTrain("Chu", chu1, sts[1].Location, sts[1], lines[1], &g)
	train2 := model.NewTrain("Cha", chu4, sts[0].Location, sts[0], lines[0], &g)
	train3 := model.NewTrain("Che", chu1, sts[3].Location, sts[3], lines[3], &g)
	train4 := model.NewTrain("Chi", chu1, sts[11].Location, sts[11], lines[0], &g)
	train5 := model.NewTrain("Cho", chu4, sts[7].Location, sts[7], lines[2], &g)
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
	for range tickerMap.C {
		go PrintMap(800, 600, sts, trains)
	}

	// Starting the server for The New Metro Times.
	ReporterServer()
}
