package main

import (
	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/models"
)

func DumpTrainsData(trains []models.Train) {
	bs := baso.NewBaso()
	for _, train := range trains {
		bs.UpdateTrain(train.Name, train.Position.X, train.Position.Y, 0.0, train.Current.ID, train.Next.ID)
	}
}
