package data

import (
	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/models"
)

func DumpTrainsData(trains []models.Train) {
	bs := baso.NewBaso()
	for _, train := range trains {
		if train.Next != nil {
			bs.UpdateTrain(train.Name, train.Position.X, train.Position.Y, 0.0, train.Current.ID, train.Next.ID)
		} else {
			bs.UpdateTrainNoNext(train.Name, train.Position.X, train.Position.Y, 0.0, train.Current.ID)
		}
	}
}

// DumpPassengersData syncs all active passengers to the database
func DumpPassengersData(stations []*models.Station, trains []models.Train) {
	bs := baso.NewBaso()

	// Collect all active passengers from stations and trains
	var allPassengers []*models.Passenger

	// Get passengers from stations (waiting)
	for _, station := range stations {
		passengers := station.GetWaitingPassengers()
		allPassengers = append(allPassengers, passengers...)
	}

	// Get passengers from trains (riding)
	for _, train := range trains {
		passengers := train.GetPassengers()
		allPassengers = append(allPassengers, passengers...)
	}

	// Batch sync to database
	err := bs.SyncPassengers(allPassengers)
	if err != nil {
		// Log error but don't crash - DB sync is non-critical
		// control.Log() would be ideal here but we don't want to spam logs
		_ = err
	}
}
