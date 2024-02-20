package main

import "internal/model"

func GenerateTestData() ([]model.Station, []model.Line) {
	sts := []model.Station{
		model.NewStation("station-1", "Station 1", model.NewVector(50.0, 350.0)),
		model.NewStation("station-2", "Station 2", model.NewVector(250.0, 200.0)),
		model.NewStation("station-3", "Station 3", model.NewVector(150.0, 100.0)),
		model.NewStation("station-4", "Station 4", model.NewVector(500.0, 50.0)),
		model.NewStation("station-5", "Station 5", model.NewVector(650.0, 150.0)),
		model.NewStation("station-6", "Station 6", model.NewVector(200.0, 400.0)),
		model.NewStation("station-7", "Station 7", model.NewVector(200.0, 500.0)),
		model.NewStation("station-8", "Station 8", model.NewVector(400.0, 450.0)),
		model.NewStation("station-9", "Station 9", model.NewVector(500.0, 350.0)),
		model.NewStation("station-10", "Station 10", model.NewVector(650.0, 300.0)),
		model.NewStation("station-11", "Station 11", model.NewVector(450.0, 150.0)),
		model.NewStation("station-12", "Station 12", model.NewVector(700.0, 50.0)),
	}
	lines := []model.Line{
		{
			Name:     "Linea 1",
			Stations: []model.Station{sts[0], sts[1], sts[3], sts[11]},
		},
		{
			Name:     "Linea 2",
			Stations: []model.Station{sts[2], sts[1], sts[5], sts[6]},
		},
		{
			Name:     "Linea 3",
			Stations: []model.Station{sts[7], sts[8], sts[9]},
		},
		{
			Name:     "Linea 4",
			Stations: []model.Station{sts[10], sts[3], sts[4]},
		},
	}
	return sts, lines
}
