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

func (bs *Baso) WipeData() error {
	tx, err := bs.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := bs.queries.WithTx(tx)

	err = qtx.DeleteAllStationLines(bs.ctx)
	if err != nil {
		return err
	}
	err = qtx.DeleteAllEdgePoints(bs.ctx)
	if err != nil {
		return err
	}
	err = qtx.DeleteAllEdges(bs.ctx)
	if err != nil {
		return err
	}
	err = qtx.DeleteAllMakes(bs.ctx)
	if err != nil {
		return err
	}
	err = qtx.DeleteAllTrains(bs.ctx)
	if err != nil {
		return err
	}
	err = qtx.DeleteAllStations(bs.ctx)
	if err != nil {
		return err
	}
	err = qtx.DeleteAllMakes(bs.ctx)
	if err != nil {
		return err
	}
	err = qtx.DeleteAllLines(bs.ctx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
