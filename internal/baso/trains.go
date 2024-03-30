package baso

import (
	"log"
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
