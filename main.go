package main

import (
	// "bufio"

	"internal/data"
	"internal/model"

	// "os"
	"time"
)

var stationHashFunction = func(station model.Station) string {
	return station.ID
}

func main() {
	// Timing and configuration
	// scnr := bufio.NewScanner(os.Stdin)
	ticker := time.NewTicker(100 * time.Millisecond)
	quit := make(chan struct{})

	// Filling graph data.
	g := data.NewGraph[model.Station](stationHashFunction)
	sts := GenerateStationsData()
	g.InsertVertices(sts)
	g.InsertEdge(sts[0], sts[3], sts[0].Location.Dist(sts[3].Location))
	g.InsertEdge(sts[0], sts[4], sts[0].Location.Dist(sts[4].Location))
	g.InsertEdge(sts[1], sts[2], sts[1].Location.Dist(sts[2].Location))
	g.InsertEdge(sts[2], sts[3], sts[2].Location.Dist(sts[3].Location))
	g.InsertEdge(sts[2], sts[5], sts[2].Location.Dist(sts[5].Location))
	g.InsertEdge(sts[2], sts[8], sts[2].Location.Dist(sts[8].Location))
	g.InsertEdge(sts[3], sts[5], sts[3].Location.Dist(sts[5].Location))
	g.InsertEdge(sts[4], sts[5], sts[4].Location.Dist(sts[5].Location))
	g.InsertEdge(sts[4], sts[9], sts[4].Location.Dist(sts[9].Location))
	g.InsertEdge(sts[5], sts[6], sts[5].Location.Dist(sts[6].Location))
	g.InsertEdge(sts[5], sts[7], sts[5].Location.Dist(sts[7].Location))

	// Creating the train and queing some destinations.
	make := model.NewMake("4-Legged-chu", "A type of fast train.", 0.01, 4)
	train := model.NewTrain("Chu", make, sts[0].Location, sts[0], &g)
	train.AddDestination(sts[2])
	train.AddDestination(sts[7])
	train.AddDestination(sts[9])

	go func() {
		for {
			select {
			case <-ticker.C:
				train.Update()
				// fmt.Println(train.Position.X, train.Position.Y)
			case <-quit:
				ticker.Stop()
				return
			}
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
