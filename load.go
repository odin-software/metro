package main

import (
	baso "github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/broadcast"
	"github.com/odin-software/metro/internal/models"
)

func LoadStations(arrs broadcast.BroadcastServer[broadcast.ADMessage[models.Train]], deps broadcast.BroadcastServer[broadcast.ADMessage[models.Train]]) []*models.Station {
	db := baso.NewBaso()
	stations := db.ListStations()
	result := make([]*models.Station, len(stations))
	for _, station := range stations {
		result = append(
			result,
			models.NewStation(station.ID, station.Name, station.Position, arrs.Subscribe(), deps.Subscribe()),
		)
	}
	return result
}

func LoadLines(arrs broadcast.BroadcastServer[broadcast.ADMessage[models.Train]], deps broadcast.BroadcastServer[broadcast.ADMessage[models.Train]]) []models.Line {
	db := baso.NewBaso()
	lines := db.ListLines()
	result := make([]models.Station, len(lines))
	for _, line := range lines {
		result = append(
			result,
			models.NewStation(station.ID, station.Name, station.Position, arrs.Subscribe(), deps.Subscribe()),
		)
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
