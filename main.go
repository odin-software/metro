package main

import (
	"fmt"
	"internal/data"
	"internal/model"
)

var stationHashFunction = func(station model.Station) string {
	return station.ID
}

func main() {
	fmt.Println("It works.")
	g := data.NewGraph[model.Station](stationHashFunction)
}
