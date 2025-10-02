package tenjin

import (
	"time"

	"github.com/odin-software/metro/internal/models"
)

// EventType represents the type of event
type EventType string

const (
	EventTypeTrainArrival   EventType = "train_arrival"
	EventTypeTrainDeparture EventType = "train_departure"
	EventTypeTrainTick      EventType = "train_tick"
	EventTypeTrainError     EventType = "train_error"
)

// Event is the base event interface that all events must implement
type Event interface {
	Type() EventType
	Timestamp() time.Time
	TrainName() string
}

// TrainArrivalEvent is emitted when a train arrives at a station
type TrainArrivalEvent struct {
	Train       string
	StationID   int64
	StationName string
	Time        time.Time
	Position    models.Vector
}

func (e TrainArrivalEvent) Type() EventType      { return EventTypeTrainArrival }
func (e TrainArrivalEvent) Timestamp() time.Time { return e.Time }
func (e TrainArrivalEvent) TrainName() string    { return e.Train }

// TrainDepartureEvent is emitted when a train departs from a station
type TrainDepartureEvent struct {
	Train       string
	StationID   int64
	StationName string
	NextStation string
	Time        time.Time
	Position    models.Vector
}

func (e TrainDepartureEvent) Type() EventType      { return EventTypeTrainDeparture }
func (e TrainDepartureEvent) Timestamp() time.Time { return e.Time }
func (e TrainDepartureEvent) TrainName() string    { return e.Train }

// TrainTickEvent is emitted periodically (every 60 ticks = 1 second) with train state
type TrainTickEvent struct {
	Train          string
	Position       models.Vector
	Velocity       models.Vector
	Speed          float64
	CurrentStation int64
	NextStation    int64 // 0 if none
	Time           time.Time
}

func (e TrainTickEvent) Type() EventType      { return EventTypeTrainTick }
func (e TrainTickEvent) Timestamp() time.Time { return e.Time }
func (e TrainTickEvent) TrainName() string    { return e.Train }

// TrainErrorEvent is emitted when a train encounters an error
type TrainErrorEvent struct {
	Train   string
	Error   string
	Context string // Additional context about what was happening
	Time    time.Time
}

func (e TrainErrorEvent) Type() EventType      { return EventTypeTrainError }
func (e TrainErrorEvent) Timestamp() time.Time { return e.Time }
func (e TrainErrorEvent) TrainName() string    { return e.Train }
