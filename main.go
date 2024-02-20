package main

import (
	// "bufio"

	"internal/model"

	// "os"
	"time"

	"github.com/VividCortex/multitick"
)

var stationHashFunction = func(station model.Station) string {
	return station.ID
}

func main() {
	// Timing and configuration
	// scnr := bufio.NewScanner(os.Stdin)
	// fmt.Println(trains[idx].Position.X, trains[idx].Position.Y)
	tick := multitick.NewTicker(15*time.Millisecond, -1*time.Millisecond)
	// ticker := time.NewTicker(15 * time.Millisecond)
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
	trainMake := model.NewMake("4-Legged-chu", "A type of fast train.", 0.01, 4)
	trains := make([]model.Train, 0)
	train := model.NewTrain("Chu", trainMake, sts[1].Location, sts[1], lines[1], &g)
	train2 := model.NewTrain("Cha", trainMake, sts[0].Location, sts[0], lines[0], &g)
	train3 := model.NewTrain("Che", trainMake, sts[3].Location, sts[3], lines[3], &g)
	train4 := model.NewTrain("Chi", trainMake, sts[11].Location, sts[11], lines[0], &g)
	train5 := model.NewTrain("Cho", trainMake, sts[7].Location, sts[7], lines[2], &g)
	trains = append(trains, train, train2, train3, train4, train5)
	// trains = append(trains, train2)

	for i := 0; i < len(trains); i++ {
		go func(idx int) {
			sub := tick.Subscribe()
			for range sub {
				trains[idx].Update()
			}
		}(i)
	}

	go func() {
		// var i int
		for range tickerMap.C {
			// i++
			// fmt.Println(i)
			PrintMap(800, 600, sts, trains)
		}
	}()
	// Adding the reporter server.
	ReporterServer()

	// for {
	// 	// getting input
	// 	fmt.Print("metro > ")
	// 	scnr.Scan()
	// 	if scnr.Text() == "stop" {
	// 		quit <- struct{}{}
	// 	}
	// 	if scnr.Text() == "exit" {
	// 		os.Exit(0)
	// 	}
	// }
}
