package baso

import (
	"database/sql"
	"log"

	"github.com/odin-software/metro/internal/dbstore"
)

type TrainsWithIds struct {
	Name             string
	X                float64
	Y                float64
	Z                float64
	CurrentStationId int64
	LineName         string
	MakeName         string
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
		Currentid: sql.NullInt64{Int64: cid},
		Nextid:    sql.NullInt64{Int64: nid},
	})
	if err != nil {
		log.Fatal(err)
	}
	return train
}
