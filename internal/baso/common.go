package baso

import (
	"log"

	"github.com/odin-software/metro/internal/models"
)

func (bs *Baso) DumpData() []models.Station {
	stations, err := bs.queries.ListStations(bs.ctx)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]models.Station, len(stations))
	for _, station := range stations {
		result = append(result, models.Station{
			ID:       station.ID,
			Name:     station.Name,
			Position: models.NewVector(station.X.Float64, station.Y.Float64),
		})
	}
	return result
}
