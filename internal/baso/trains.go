package baso

import (
	"database/sql"
	"log"

	"github.com/odin-software/metro/internal/dbstore"
)

type TrainsWithIds struct {
	ID               int64   `json:"id"`
	Name             string  `json:"name"`
	X                float64 `json:"x"`
	Y                float64 `json:"y"`
	Z                float64 `json:"z"`
	CurrentStationId int64   `json:"currentId"`
	LineName         string  `json:"line"`
	MakeName         string  `json:"make"`
}

func (bs *Baso) ListTrainsFull() []TrainsWithIds {
	trains, err := bs.queries.GetAllTrainsFull(bs.ctx)
	if err != nil {
		log.Fatal(err)
	}
	result := make([]TrainsWithIds, 0)
	for _, train := range trains {
		result = append(
			result,
			TrainsWithIds{
				ID:               train.ID,
				Name:             train.Name,
				X:                train.X,
				Y:                train.Y,
				Z:                train.Z,
				CurrentStationId: train.Stationid,
				LineName:         train.Linename,
				MakeName:         train.Makename,
			},
		)
	}
	return result
}

func (bs *Baso) UpdateTrain(name string, x float64, y float64, z float64, cid int64, nid int64) int64 {
	train, err := bs.queries.UpdateTrain(bs.ctx, dbstore.UpdateTrainParams{
		Name:      name,
		Name_2:    name,
		X:         x,
		Y:         y,
		Z:         z,
		Currentid: sql.NullInt64{Int64: cid, Valid: true},
		Nextid:    sql.NullInt64{Int64: nid, Valid: true},
	})
	if err != nil {
		log.Fatal(err)
	}
	return train
}

func (bs *Baso) UpdateTrainNoNext(name string, x float64, y float64, z float64, cid int64) int64 {
	train, err := bs.queries.UpdateTrain(bs.ctx, dbstore.UpdateTrainParams{
		Name:      name,
		Name_2:    name,
		X:         x,
		Y:         y,
		Z:         z,
		Currentid: sql.NullInt64{Int64: cid, Valid: true},
	})
	if err != nil {
		log.Fatal(err)
	}
	return train
}

func (bs *Baso) MoveTrainToLine(trainId, lineId int64) error {
	err := bs.queries.ChangeTrainToLine(bs.ctx, dbstore.ChangeTrainToLineParams{
		ID: trainId,
		Lineid: sql.NullInt64{
			Int64: lineId,
			Valid: true,
		},
	})
	if err != nil {
		return err
	}
	return nil
}
