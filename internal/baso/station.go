package baso

import (
	"log"

	"github.com/odin-software/metro/internal/models"
)

func (bs *Baso) ListStations() []models.Station {
	stations, err := bs.queries.ListStations(bs.ctx)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]models.Station, 0)
	for _, station := range stations {
		result = append(result, models.Station{
			ID:       station.ID,
			Name:     station.Name,
			Position: models.NewVector(station.X.Float64, station.Y.Float64),
		})
	}
	return result
}

func (bs *Baso) GetStationById(id int64) models.Station {
	station, err := bs.queries.GetStationById(bs.ctx, id)
	if err != nil {
		log.Fatal(err)
	}
	return models.Station{
		ID:       station.ID,
		Name:     station.Name,
		Position: models.NewVector(station.X.Float64, station.Y.Float64),
	}
}
