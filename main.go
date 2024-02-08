package main

import (
	"data"
	"fmt"
	"model"
)

var stationHashFunction = func(station model.Station) string {
	return station.ID
}

func main() {
	// Filling graph data.
	g := data.NewGraph[model.Station](stationHashFunction)
	sts := GenerateStationsData()
	g.InsertVertices(sts)

	make := model.NewMake("4-Legged", "A type of fast train.", 0.01, 4)
	train := model.NewTrain("Chu", make, sts[0].Location, &g)
	train.AddDestination(sts[2])
	train.AddDestination(sts[7])
	train.AddDestination(sts[9])

	for {
		train.Update()
		fmt.Println(train.Position.X, train.Position.Y)
	}
}
