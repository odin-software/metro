package display

import (
	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/dbstore"
)

// BasoScheduleAdapter adapts baso for the display package
type BasoScheduleAdapter struct {
	baso *baso.Baso
}

// NewBasoScheduleAdapter creates a new adapter
func NewBasoScheduleAdapter() *BasoScheduleAdapter {
	return &BasoScheduleAdapter{
		baso: baso.NewBaso(),
	}
}

// GetScheduleByTrainAndStation looks up a specific schedule entry
func (a *BasoScheduleAdapter) GetScheduleByTrainAndStation(trainID, stationID int64) (Schedule, error) {
	dbSched, err := a.baso.GetScheduleByTrainAndStation(trainID, stationID)
	if err != nil {
		return Schedule{}, err
	}
	return dbstoreToDisplaySchedule(dbSched), nil
}

// GetScheduleForTrain gets all scheduled stops for a train
func (a *BasoScheduleAdapter) GetScheduleForTrain(trainID int64) ([]Schedule, error) {
	dbScheds, err := a.baso.GetScheduleForTrain(trainID)
	if err != nil {
		return nil, err
	}

	result := make([]Schedule, len(dbScheds))
	for i, s := range dbScheds {
		result[i] = dbstoreToDisplaySchedule(s)
	}
	return result, nil
}

// GetScheduleForStation gets all scheduled arrivals at a station
func (a *BasoScheduleAdapter) GetScheduleForStation(stationID int64) ([]Schedule, error) {
	dbScheds, err := a.baso.GetScheduleForStation(stationID)
	if err != nil {
		return nil, err
	}

	result := make([]Schedule, len(dbScheds))
	for i, s := range dbScheds {
		result[i] = dbstoreToDisplaySchedule(s)
	}
	return result, nil
}

// dbstoreToDisplaySchedule converts dbstore.Schedule to display.Schedule
func dbstoreToDisplaySchedule(s dbstore.Schedule) Schedule {
	return Schedule{
		TrainID:       s.TrainID,
		StationID:     s.StationID,
		ScheduledTime: int(s.ScheduledTime),
		SequenceOrder: s.SequenceOrder,
	}
}
