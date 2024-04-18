package baso

import (
	"database/sql"
	"log"
	"reflect"

	models "github.com/odin-software/metro/internal/models"
)

type LineWithEdges struct {
	Name   string          `json:"name"`
	Points []models.Vector `json:"points"`
}

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

func (bs *Baso) ListLinesWithPoints() ([]LineWithEdges, error) {
	lines, err := bs.queries.ListLines(bs.ctx)
	if err != nil {
		return nil, err
	}
	result := make([]LineWithEdges, 0)
	for _, line := range lines {
		points, err := bs.queries.GetPointsFromLine(bs.ctx, line.ID)
		if err != nil {
			return nil, err
		}
		vectors := make([]models.Vector, 0)
		for _, p := range points {
			vx := reflect.ValueOf(p.X)
			x := vx.Convert(reflect.TypeOf(float64(0))).Float()
			vy := reflect.ValueOf(p.Y)
			y := vy.Convert(reflect.TypeOf(float64(0))).Float()
			vectors = append(vectors, models.NewVector(x, y))
		}
		result = append(result, LineWithEdges{
			Name:   line.Name,
			Points: vectors,
		})
	}
	return result, nil
}
