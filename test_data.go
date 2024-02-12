package main

import "internal/model"

func GenerateStationsData() []model.Station {
	sts := []model.Station{
		model.NewStation("station-1", "Station 1", model.NewVector(100.0, 250.0)),
		model.NewStation("station-2", "Station 2", model.NewVector(150.0, 600.0)),
		model.NewStation("station-3", "Station 3", model.NewVector(300.0, 550.0)),
		model.NewStation("station-4", "Station 4", model.NewVector(300.0, 300.0)),
		model.NewStation("station-5", "Station 5", model.NewVector(250.0, 150.0)),
		model.NewStation("station-6", "Station 6", model.NewVector(450.0, 300.0)),
		model.NewStation("station-7", "Station 7", model.NewVector(550.0, 550.0)),
		model.NewStation("station-8", "Station 8", model.NewVector(700.0, 250.0)),
		model.NewStation("station-9", "Station 9", model.NewVector(400.0, 700.0)),
		model.NewStation("station-10", "Station 10", model.NewVector(500.0, 100.0)),
	}
	return sts
}
