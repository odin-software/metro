package main

import (
	"github.com/odin-software/metro/internal/broadcast"
	"github.com/odin-software/metro/internal/models"
)

func GenerateTestData(arrs broadcast.BroadcastServer[broadcast.ADMessage[models.Train]], deps broadcast.BroadcastServer[broadcast.ADMessage[models.Train]]) ([]*models.Station, []models.Line) {
	sts := []*models.Station{
		models.NewStation(1, "Station 1", models.NewVector(50.0, 350.0), arrs.Subscribe(), deps.Subscribe()),
		models.NewStation(2, "Station 2", models.NewVector(250.0, 200.0), arrs.Subscribe(), deps.Subscribe()),
		models.NewStation(3, "Station 3", models.NewVector(150.0, 100.0), arrs.Subscribe(), deps.Subscribe()),
		models.NewStation(4, "Station 4", models.NewVector(500.0, 50.0), arrs.Subscribe(), deps.Subscribe()),
		models.NewStation(5, "Station 5", models.NewVector(650.0, 150.0), arrs.Subscribe(), deps.Subscribe()),
		models.NewStation(6, "Station 6", models.NewVector(200.0, 400.0), arrs.Subscribe(), deps.Subscribe()),
		models.NewStation(7, "Station 7", models.NewVector(200.0, 500.0), arrs.Subscribe(), deps.Subscribe()),
		models.NewStation(8, "Station 8", models.NewVector(400.0, 450.0), arrs.Subscribe(), deps.Subscribe()),
		models.NewStation(9, "Station 9", models.NewVector(500.0, 350.0), arrs.Subscribe(), deps.Subscribe()),
		models.NewStation(10, "Station 10", models.NewVector(650.0, 300.0), arrs.Subscribe(), deps.Subscribe()),
		models.NewStation(11, "Station 11", models.NewVector(450.0, 150.0), arrs.Subscribe(), deps.Subscribe()),
		models.NewStation(12, "Station 12", models.NewVector(700.0, 50.0), arrs.Subscribe(), deps.Subscribe()),
	}
	lines := []models.Line{
		{
			Name:     "Linea 1",
			Stations: []models.Station{*sts[0], *sts[1], *sts[3], *sts[11]},
		},
		{
			Name:     "Linea 2",
			Stations: []models.Station{*sts[2], *sts[1], *sts[5], *sts[6]},
		},
		{
			Name:     "Linea 3",
			Stations: []models.Station{*sts[7], *sts[8], *sts[9]},
		},
		{
			Name:     "Linea 4",
			Stations: []models.Station{*sts[10], *sts[3], *sts[4]},
		},
	}
	return sts, lines
}
