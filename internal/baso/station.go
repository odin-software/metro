package baso

import (
	"database/sql"

	"github.com/odin-software/metro/internal/dbstore"
	"github.com/odin-software/metro/internal/models"
)

type GetStation struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	Color    string        `json:"color"`
	Position models.Vector `json:"position"`
}

type CreateStation struct {
	Name  string  `json:"name"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Z     float64 `json:"z"`
	Color string  `json:"color"`
}

func (bs *Baso) ListStations() ([]GetStation, error) {
	stations, err := bs.queries.ListStations(bs.ctx)
	if err != nil {
		return nil, err
	}
	result := make([]GetStation, 0)
	for _, station := range stations {
		result = append(result, GetStation{
			ID:       station.ID,
			Name:     station.Name,
			Color:    station.Color.String,
			Position: models.NewVector(station.X.Float64, station.Y.Float64),
		})
	}
	return result, nil
}

func (bs *Baso) GetStationById(id int64) (GetStation, error) {
	station, err := bs.queries.GetStationById(bs.ctx, id)
	if err != nil {
		return GetStation{}, err
	}
	return GetStation{
		ID:       station.ID,
		Name:     station.Name,
		Color:    station.Color.String,
		Position: models.NewVector(station.X.Float64, station.Y.Float64),
	}, err
}

func (bs *Baso) CreateStation(name, color string, x, y, z float64) error {
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
		Color: sql.NullString{
			String: color,
			Valid:  true,
		},
	})

	return err
}

func (bs *Baso) CreateStations(sts []CreateStation) ([]GetStation, error) {
	tx, err := bs.db.Begin()
	if err != nil {
		return nil, err
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
			Color: sql.NullString{
				String: st.Color,
				Valid:  true,
			},
		})
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	newStations, err := bs.ListStations()
	if err != nil {
		return nil, err
	}

	return newStations, nil
}
