package baso

import (
	"database/sql"
	"log"

	"github.com/odin-software/metro/internal/dbstore"
	"github.com/odin-software/metro/internal/models"
)

type CreateStation struct {
	Name string  `json:"name"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
}

func (bs *Baso) ListStations() ([]models.Station, error) {
	stations, err := bs.queries.ListStations(bs.ctx)
	if err != nil {
		return nil, err
	}
	result := make([]models.Station, 0)
	for _, station := range stations {
		result = append(result, models.Station{
			ID:       station.ID,
			Name:     station.Name,
			Position: models.NewVector(station.X.Float64, station.Y.Float64),
		})
	}
	return result, nil
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

func (bs *Baso) CreateStation(name string, x, y, z float64) error {
	_, err := bs.queries.CreateStation(bs.ctx, dbstore.CreateStationParams{
		Name: name,
		X: sql.NullFloat64{
			Float64: x,
			Valid:   true,
		},
		Y: sql.NullFloat64{
			Float64: y,
			Valid:   true,
		},
		Z: sql.NullFloat64{
			Float64: 0,
			Valid:   true,
		},
	})

	return err
}

func (bs *Baso) CreateStations(sts []CreateStation) error {
	tx, err := bs.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := bs.queries.WithTx(tx)

	for _, st := range sts {
		_, err := qtx.CreateStation(bs.ctx, dbstore.CreateStationParams{
			Name: st.Name,
			X: sql.NullFloat64{
				Float64: st.X,
				Valid:   true,
			},
			Y: sql.NullFloat64{
				Float64: st.Y,
				Valid:   true,
			},
			Z: sql.NullFloat64{
				Float64: 0,
				Valid:   true,
			},
		})
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
