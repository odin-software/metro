package analysis

import (
	"github.com/odin-software/metro/internal/baso"
)

// BasoScheduleAdapter adapts baso.Baso to implement the ScheduleDB interface
type BasoScheduleAdapter struct {
	baso *baso.Baso
}

// NewBasoScheduleAdapter creates a new adapter
func NewBasoScheduleAdapter() *BasoScheduleAdapter {
	return &BasoScheduleAdapter{
		baso: baso.NewBaso(),
	}
}

// GetScheduleByTrainAndStation looks up a schedule entry
func (a *BasoScheduleAdapter) GetScheduleByTrainAndStation(trainID, stationID int64) (Schedule, error) {
	dbSchedule, err := a.baso.GetScheduleByTrainAndStation(trainID, stationID)
	if err != nil {
		return Schedule{}, err
	}

	// Convert dbstore.Schedule to analysis.Schedule
	return Schedule{
		TrainID:       dbSchedule.TrainID,
		StationID:     dbSchedule.StationID,
		ScheduledTime: int(dbSchedule.ScheduledTime),
	}, nil
}
