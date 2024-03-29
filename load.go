package main

import (
	baso "github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/broadcast"
	"github.com/odin-software/metro/internal/models"
)

func LoadStations(arrs broadcast.BroadcastServer[broadcast.ADMessage[models.Train]], deps broadcast.BroadcastServer[broadcast.ADMessage[models.Train]]) []*models.Station {
	db := baso.NewBaso()
	stations := db.ListStations()
	result := make([]*models.Station, 0)
	for _, station := range stations {
		result = append(
			result,
			models.NewStation(station.ID, station.Name, station.Position, arrs.Subscribe(), deps.Subscribe()),
		)
	}
	return result
}

func LoadLines() []models.Line {
	db := baso.NewBaso()
	lines := db.ListLinesWithStations()
	result := make([]models.Line, len(lines))
	for _, line := range lines {
		result = append(
			result,
			models.Line{
				Name:     line.Name,
				Stations: line.Stations,
			},
		)
	}

	return result
}
