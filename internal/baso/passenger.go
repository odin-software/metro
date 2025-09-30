package baso

import (
	"database/sql"

	"github.com/odin-software/metro/internal/dbstore"
	"github.com/odin-software/metro/internal/models"
)

// SavePassenger persists a passenger to the database
func (bs *Baso) SavePassenger(passenger *models.Passenger) error {
	// Note: current_train_id is set to NULL for now
	// Train struct doesn't expose ID, only Name

	_, err := bs.queries.CreatePassenger(bs.ctx, dbstore.CreatePassengerParams{
		ID:                   passenger.ID,
		Name:                 passenger.Name,
		CurrentStationID:     passenger.CurrentStation.ID,
		DestinationStationID: passenger.DestinationStation.ID,
		State:                string(passenger.State),
		Sentiment:            passenger.Sentiment,
		SpawnTime:            passenger.WaitStartTime,
	})

	return err
}

// UpdatePassengerState updates passenger state and sentiment
func (bs *Baso) UpdatePassengerState(passengerID string, state models.PassengerState, sentiment float64) error {
	return bs.queries.UpdatePassengerState(bs.ctx, dbstore.UpdatePassengerStateParams{
		ID:        passengerID,
		State:     string(state),
		Sentiment: sentiment,
	})
}

// UpdatePassengerBoarding marks a passenger as boarding a train
// Note: trainID set to NULL since Train struct doesn't expose ID
func (bs *Baso) UpdatePassengerBoarding(passengerID string) error {
	return bs.queries.UpdatePassengerBoarding(bs.ctx, dbstore.UpdatePassengerBoardingParams{
		ID:             passengerID,
		CurrentTrainID: sql.NullInt64{Valid: false},
	})
}

// UpdatePassengerDisembarking marks a passenger as disembarking
func (bs *Baso) UpdatePassengerDisembarking(passengerID string, stationID int64, state models.PassengerState) error {
	return bs.queries.UpdatePassengerDisembarking(bs.ctx, dbstore.UpdatePassengerDisembarkingParams{
		ID:               passengerID,
		State:            string(state),
		CurrentStationID: stationID,
	})
}

// DeletePassenger removes a passenger from the database
func (bs *Baso) DeletePassenger(passengerID string) error {
	return bs.queries.DeletePassenger(bs.ctx, passengerID)
}

// GetActivePassengers retrieves all active passengers
func (bs *Baso) GetActivePassengers() ([]dbstore.Passenger, error) {
	return bs.queries.GetAllActivePassengers(bs.ctx)
}

// CountActivePassengers returns the number of active passengers
func (bs *Baso) CountActivePassengers() (int64, error) {
	return bs.queries.CountActivePassengers(bs.ctx)
}

// CreatePassengerEvent logs an event for a passenger
func (bs *Baso) CreatePassengerEvent(passengerID, eventType string, stationID *int64, trainID *int64, sentiment float64, metadata string) error {
	params := dbstore.CreatePassengerEventParams{
		PassengerID: passengerID,
		EventType:   eventType,
		Sentiment:   sql.NullFloat64{Float64: sentiment, Valid: true},
		Metadata:    sql.NullString{String: metadata, Valid: metadata != ""},
	}

	if stationID != nil {
		params.StationID = sql.NullInt64{Int64: *stationID, Valid: true}
	}

	if trainID != nil {
		params.TrainID = sql.NullInt64{Int64: *trainID, Valid: true}
	}

	_, err := bs.queries.CreatePassengerEvent(bs.ctx, params)
	return err
}

// GetPassengerEvents retrieves all events for a passenger
func (bs *Baso) GetPassengerEvents(passengerID string) ([]dbstore.PassengerEvent, error) {
	return bs.queries.GetPassengerEvents(bs.ctx, passengerID)
}

// DeletePassengerEvents removes all events for a passenger
func (bs *Baso) DeletePassengerEvents(passengerID string) error {
	return bs.queries.DeletePassengerEvents(bs.ctx, passengerID)
}

// SyncPassengers batch syncs multiple passengers to database
func (bs *Baso) SyncPassengers(passengers []*models.Passenger) error {
	tx, err := bs.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := bs.queries.WithTx(tx)

	// Delete all existing passengers first (since we're syncing state)
	err = qtx.DeleteAllPassengers(bs.ctx)
	if err != nil {
		return err
	}

	// Insert current active passengers
	for _, p := range passengers {
		if p.State == models.PassengerStateArrived {
			continue // Skip arrived passengers
		}

		// Note: current_train_id set to NULL (Train struct doesn't expose ID)

		_, err := qtx.CreatePassenger(bs.ctx, dbstore.CreatePassengerParams{
			ID:                   p.ID,
			Name:                 p.Name,
			CurrentStationID:     p.CurrentStation.ID,
			DestinationStationID: p.DestinationStation.ID,
			State:                string(p.State),
			Sentiment:            p.Sentiment,
			SpawnTime:            p.WaitStartTime,
		})

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// LogPassengerEvent is a helper to log events from passenger lifecycle
func (bs *Baso) LogPassengerEvent(event interface{}) error {
	// Parse event type and extract fields
	switch e := event.(type) {
	case map[string]interface{}:
		eventType, ok := e["Type"].(string)
		if !ok {
			return nil // Skip invalid events
		}

		passengerID, ok := e["PassengerID"].(string)
		if !ok {
			passengerID, ok = e["PassengerName"].(string) // Fallback
			if !ok {
				return nil
			}
		}

		var stationID *int64
		if sid, ok := e["StationID"].(int64); ok {
			stationID = &sid
		}

		var trainID *int64
		var trainName string
		if e["TrainName"] != nil {
			trainName = e["TrainName"].(string)
		}

		sentiment := 0.0
		if s, ok := e["Sentiment"].(float64); ok {
			sentiment = s
		}

		// Store train name in metadata if available
		metadata := ""
		if trainName != "" {
			metadata = trainName
		}

		return bs.CreatePassengerEvent(passengerID, eventType, stationID, trainID, sentiment, metadata)
	}

	return nil
}
