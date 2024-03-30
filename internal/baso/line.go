package baso

import (
	"database/sql"
	"log"

	models "github.com/odin-software/metro/internal/models"
)

func (bs *Baso) ListLinesWithStations() []models.Line {
	lines, err := bs.queries.ListLines(bs.ctx)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]models.Line, 0)
	for _, line := range lines {
		stationsInLine, err := bs.queries.GetStationsFromLine(bs.ctx, sql.NullInt64{Int64: line.ID, Valid: true})
		if err != nil {
			log.Fatal(err)
		}
		stations := make([]models.Station, 0)
		for _, station := range stationsInLine {
			stations = append(stations, models.Station{
				ID:       station.ID,
				Name:     station.Name,
				Position: models.NewVector(station.X.Float64, station.Y.Float64),
			})
		}
		result = append(result, models.Line{
			Name:     line.Name,
			Stations: stations,
		})
	}
	return result
}
