package models

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// PassengerState represents the current state of a passenger
type PassengerState string

const (
	PassengerStateWaiting      PassengerState = "waiting"      // Waiting at station
	PassengerStateBoarding     PassengerState = "boarding"     // In process of boarding train
	PassengerStateRiding       PassengerState = "riding"       // On the train
	PassengerStateDisembarking PassengerState = "disembarking" // In process of leaving train
	PassengerStateArrived      PassengerState = "arrived"      // Reached destination
)

// Passenger represents a person using the transit system
type Passenger struct {
	ID                 string
	Name               string
	Position           Vector
	CurrentStation     *Station
	DestinationStation *Station
	CurrentTrain       *Train  // nil if not on a train
	Sentiment          float64 // 0-100, higher is better
	State              PassengerState
	WaitStartTime      time.Time          // When they started waiting
	JourneyStartTime   time.Time          // When they spawned/started journey
	lastSentimentDrop  time.Time          // Last time sentiment was decreased
	eventChannel       chan<- interface{} // Channel to send events to Tenjin
	Drawing                               // For future visualization
}

// NewPassenger creates a new passenger
func NewPassenger(
	id string,
	name string,
	currentStation *Station,
	destinationStation *Station,
	eventChannel chan<- interface{},
) *Passenger {
	p := &Passenger{
		ID:                 id,
		Name:               name,
		Position:           currentStation.Position, // Start at station position
		CurrentStation:     currentStation,
		DestinationStation: destinationStation,
		CurrentTrain:       nil,
		Sentiment:          100.0, // Start with perfect satisfaction
		State:              PassengerStateWaiting,
		WaitStartTime:      time.Now(),
		JourneyStartTime:   time.Time{}, // Will be set when boarding
		lastSentimentDrop:  time.Now(),
		eventChannel:       eventChannel,
	}

	// Emit spawn event
	p.emitSpawnEvent()

	return p
}

// UpdateSentiment adjusts passenger sentiment based on waiting/riding conditions
func (p *Passenger) UpdateSentiment(deltaTime time.Duration) {
	// Only update sentiment every 5 seconds to avoid rapid drops
	if time.Since(p.lastSentimentDrop) < 5*time.Second {
		return
	}

	switch p.State {
	case PassengerStateWaiting:
		// Lose 2 points every 5 seconds of waiting
		waitTime := time.Since(p.WaitStartTime)
		if waitTime >= 5*time.Second {
			p.Sentiment -= 2.0
			if p.Sentiment < 0 {
				p.Sentiment = 0
			}
			p.lastSentimentDrop = time.Now()

			// Emit frustration event when sentiment drops below 50
			if p.Sentiment < 50 {
				p.emitFrustrationEvent()
			}
		}

	case PassengerStateRiding:
		// Minor sentiment decrease for long journeys
		journeyTime := time.Since(p.JourneyStartTime)
		if journeyTime > 15*time.Second {
			p.Sentiment -= 0.5
			if p.Sentiment < 0 {
				p.Sentiment = 0
			}
			p.lastSentimentDrop = time.Now()

			// Extra penalty if train is crowded
			if p.CurrentTrain != nil && p.CurrentTrain.IsCrowded() {
				p.Sentiment -= 1.0
			}
		}
	}
}

// StartWaiting sets passenger to waiting state
func (p *Passenger) StartWaiting() {
	p.State = PassengerStateWaiting
	p.WaitStartTime = time.Now()
	p.lastSentimentDrop = time.Now()
	p.emitWaitEvent()
}

// BoardTrain puts passenger on a train
func (p *Passenger) BoardTrain(train *Train) {
	p.State = PassengerStateBoarding
	p.CurrentTrain = train
	p.Position = train.Position
	p.State = PassengerStateRiding
	p.JourneyStartTime = time.Now()  // Start tracking journey time
	p.WaitStartTime = time.Time{}    // Clear wait timer
	p.lastSentimentDrop = time.Now() // Reset sentiment drop timer
	p.emitBoardEvent()
}

// DisembarkTrain removes passenger from train
func (p *Passenger) DisembarkTrain(station *Station) {
	p.State = PassengerStateDisembarking
	p.CurrentTrain = nil
	p.CurrentStation = station
	p.Position = station.Position

	// Always emit disembark event when leaving a train
	p.emitDisembarkEvent()

	// Check if arrived at destination
	if station.ID == p.DestinationStation.ID {
		p.State = PassengerStateArrived
		p.emitArriveEvent()
		p.JourneyStartTime = time.Time{} // Clear journey timer
	} else {
		// Transfer - start waiting again
		p.State = PassengerStateWaiting
		p.WaitStartTime = time.Now()
		p.JourneyStartTime = time.Time{} // Reset for next leg
		p.lastSentimentDrop = time.Now() // Reset sentiment drop timer
		p.emitWaitEvent()
	}
}

// GetSentimentCategory returns a human-readable sentiment
func (p *Passenger) GetSentimentCategory() string {
	switch {
	case p.Sentiment >= 80:
		return "Happy"
	case p.Sentiment >= 60:
		return "Satisfied"
	case p.Sentiment >= 40:
		return "Neutral"
	case p.Sentiment >= 20:
		return "Frustrated"
	default:
		return "Angry"
	}
}

// String returns a string representation of the passenger
func (p *Passenger) String() string {
	return fmt.Sprintf("%s (%s) at %s → %s [%s, %.0f%%]",
		p.Name,
		p.ID,
		p.CurrentStation.Name,
		p.DestinationStation.Name,
		p.State,
		p.Sentiment,
	)
}

// Event emission methods

func (p *Passenger) emitSpawnEvent() {
	if p.eventChannel == nil {
		return
	}

	event := struct {
		Type            string
		PassengerID     string
		PassengerName   string
		StationID       int64
		StationName     string
		DestinationID   int64
		DestinationName string
		Time            time.Time
	}{
		Type:            "passenger_spawn",
		PassengerID:     p.ID,
		PassengerName:   p.Name,
		StationID:       p.CurrentStation.ID,
		StationName:     p.CurrentStation.Name,
		DestinationID:   p.DestinationStation.ID,
		DestinationName: p.DestinationStation.Name,
		Time:            time.Now(),
	}

	select {
	case p.eventChannel <- event:
	default:
		// Channel full, skip event
	}
}

func (p *Passenger) emitWaitEvent() {
	if p.eventChannel == nil {
		return
	}

	event := struct {
		Type          string
		PassengerID   string
		PassengerName string
		StationID     int64
		StationName   string
		WaitDuration  time.Duration
		Sentiment     float64
		Time          time.Time
	}{
		Type:          "passenger_wait",
		PassengerID:   p.ID,
		PassengerName: p.Name,
		StationID:     p.CurrentStation.ID,
		StationName:   p.CurrentStation.Name,
		WaitDuration:  time.Since(p.WaitStartTime),
		Sentiment:     p.Sentiment,
		Time:          time.Now(),
	}

	select {
	case p.eventChannel <- event:
	default:
		// Channel full, skip event
	}
}

func (p *Passenger) emitBoardEvent() {
	if p.eventChannel == nil {
		return
	}

	event := struct {
		Type          string
		PassengerID   string
		PassengerName string
		TrainName     string
		StationID     int64
		StationName   string
		Sentiment     float64
		Time          time.Time
	}{
		Type:          "passenger_board",
		PassengerID:   p.ID,
		PassengerName: p.Name,
		TrainName:     p.CurrentTrain.Name,
		StationID:     p.CurrentStation.ID,
		StationName:   p.CurrentStation.Name,
		Sentiment:     p.Sentiment,
		Time:          time.Now(),
	}

	select {
	case p.eventChannel <- event:
	default:
		// Channel full, skip event
	}
}

func (p *Passenger) emitDisembarkEvent() {
	if p.eventChannel == nil {
		return
	}

	event := struct {
		Type          string
		PassengerID   string
		PassengerName string
		StationID     int64
		StationName   string
		Sentiment     float64
		Time          time.Time
	}{
		Type:          "passenger_disembark",
		PassengerID:   p.ID,
		PassengerName: p.Name,
		StationID:     p.CurrentStation.ID,
		StationName:   p.CurrentStation.Name,
		Sentiment:     p.Sentiment,
		Time:          time.Now(),
	}

	select {
	case p.eventChannel <- event:
	default:
		// Channel full, skip event
	}
}

func (p *Passenger) emitArriveEvent() {
	if p.eventChannel == nil {
		return
	}

	journeyDuration := time.Since(p.JourneyStartTime)

	event := struct {
		Type            string
		PassengerID     string
		PassengerName   string
		DestinationID   int64
		DestinationName string
		JourneyDuration time.Duration
		Sentiment       float64
		Time            time.Time
	}{
		Type:            "passenger_arrive",
		PassengerID:     p.ID,
		PassengerName:   p.Name,
		DestinationID:   p.DestinationStation.ID,
		DestinationName: p.DestinationStation.Name,
		JourneyDuration: journeyDuration,
		Sentiment:       p.Sentiment,
		Time:            time.Now(),
	}

	select {
	case p.eventChannel <- event:
	default:
		// Channel full, skip event
	}
}

func (p *Passenger) emitFrustrationEvent() {
	if p.eventChannel == nil {
		return
	}

	event := struct {
		Type          string
		PassengerID   string
		PassengerName string
		Sentiment     float64
		Category      string
		Reason        string
		Time          time.Time
	}{
		Type:          "passenger_frustration",
		PassengerID:   p.ID,
		PassengerName: p.Name,
		Sentiment:     p.Sentiment,
		Category:      p.GetSentimentCategory(),
		Reason:        fmt.Sprintf("Waiting for %.0f seconds", time.Since(p.WaitStartTime).Seconds()),
		Time:          time.Now(),
	}

	select {
	case p.eventChannel <- event:
	default:
		// Channel full, skip event
	}
}

// Display methods (for future visualization)

func (p *Passenger) Update() {
	p.Drawing.Counter++
	// Update sentiment based on time
	p.UpdateSentiment(time.Second / 60) // Called every tick
}

func (p *Passenger) Draw(screen *ebiten.Image) {
	// Future: Draw passenger sprite at position
	// For now, passengers are not rendered
}
