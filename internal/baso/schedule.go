package baso

import (
	"github.com/odin-software/metro/internal/dbstore"
)

// GetScheduleForTrain retrieves all scheduled stops for a train, ordered by sequence
func (bs *Baso) GetScheduleForTrain(trainID int64) ([]dbstore.Schedule, error) {
	return bs.queries.GetScheduleForTrain(bs.ctx, trainID)
}

// GetScheduleForStation retrieves all scheduled arrivals at a station
func (bs *Baso) GetScheduleForStation(stationID int64) ([]dbstore.Schedule, error) {
	return bs.queries.GetScheduleForStation(bs.ctx, stationID)
}

// GetNextScheduledStop retrieves the next stop for a train after a given sequence order
func (bs *Baso) GetNextScheduledStop(trainID int64, currentSequence int64) (dbstore.Schedule, error) {
	return bs.queries.GetNextScheduledStop(bs.ctx, dbstore.GetNextScheduledStopParams{
		TrainID:       trainID,
		SequenceOrder: currentSequence,
	})
}

// GetScheduleByTrainAndStation finds a specific schedule entry
func (bs *Baso) GetScheduleByTrainAndStation(trainID, stationID int64) (dbstore.Schedule, error) {
	return bs.queries.GetScheduleByTrainAndStation(bs.ctx, dbstore.GetScheduleByTrainAndStationParams{
		TrainID:   trainID,
		StationID: stationID,
	})
}

// GetAllSchedules retrieves all schedules, ordered by train and sequence
func (bs *Baso) GetAllSchedules() ([]dbstore.Schedule, error) {
	return bs.queries.GetAllSchedules(bs.ctx)
}
