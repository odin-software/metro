package models

import (
	"image"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/odin-software/metro/internal/assets"
)

type Station struct {
	ID                int64        `json:"id"`
	Name              string       `json:"name"`
	Position          Vector       `json:"position"`
	WaitingPassengers []*Passenger // Passengers waiting at this station
	passengerMutex    sync.RWMutex // Thread safety for passenger operations
	Drawing
}

func NewStation(id int64, name string, location Vector) *Station {
	img, frameWidth, frameHeight, frameCount := assets.GetStationSprite()
	return &Station{
		ID:                id,
		Name:              name,
		Position:          location,
		WaitingPassengers: make([]*Passenger, 0),
		Drawing: Drawing{
			Counter:     0,
			FrameWidth:  frameWidth,
			FrameHeight: frameHeight,
			FrameCount:  frameCount,
			Sprite:      img,
		},
	}
}

// Passenger management methods

// AddPassenger adds a passenger to the station's waiting queue
func (st *Station) AddPassenger(passenger *Passenger) {
	st.passengerMutex.Lock()
	defer st.passengerMutex.Unlock()

	st.WaitingPassengers = append(st.WaitingPassengers, passenger)
	passenger.CurrentStation = st
	passenger.Position = st.Position
}

// RemovePassenger removes a passenger from the station's waiting queue
func (st *Station) RemovePassenger(passenger *Passenger) bool {
	st.passengerMutex.Lock()
	defer st.passengerMutex.Unlock()

	for i, p := range st.WaitingPassengers {
		if p.ID == passenger.ID {
			// Remove passenger by swapping with last and truncating
			st.WaitingPassengers[i] = st.WaitingPassengers[len(st.WaitingPassengers)-1]
			st.WaitingPassengers = st.WaitingPassengers[:len(st.WaitingPassengers)-1]
			return true
		}
	}
	return false
}

// GetWaitingPassengersCount returns the number of passengers waiting
func (st *Station) GetWaitingPassengersCount() int {
	st.passengerMutex.RLock()
	defer st.passengerMutex.RUnlock()
	return len(st.WaitingPassengers)
}

// GetWaitingPassengers returns a copy of the waiting passengers slice (thread-safe)
func (st *Station) GetWaitingPassengers() []*Passenger {
	st.passengerMutex.RLock()
	defer st.passengerMutex.RUnlock()

	passengers := make([]*Passenger, len(st.WaitingPassengers))
	copy(passengers, st.WaitingPassengers)
	return passengers
}

// GetPassengersForDestination returns passengers waiting to go to a specific station
func (st *Station) GetPassengersForDestination(destinationID int64) []*Passenger {
	st.passengerMutex.RLock()
	defer st.passengerMutex.RUnlock()

	result := make([]*Passenger, 0)
	for _, p := range st.WaitingPassengers {
		if p.DestinationStation.ID == destinationID {
			result = append(result, p)
		}
	}
	return result
}

// Drawing methods

func (st *Station) Update() {
	st.Drawing.Counter++

	// Update all waiting passengers
	st.passengerMutex.RLock()
	passengers := st.WaitingPassengers
	st.passengerMutex.RUnlock()

	for _, passenger := range passengers {
		if passenger.State == PassengerStateWaiting {
			passenger.Update()
		}
	}
}

func (st *Station) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(st.FrameWidth)/2, -float64(st.FrameHeight)/2)
	op.GeoM.Translate(st.Position.X, st.Position.Y)
	i := (st.Counter / st.FrameCount) % st.FrameCount
	sx, sy := 0+i*st.FrameWidth, 0
	screen.DrawImage(st.Sprite.SubImage(image.Rect(sx, sy, sx+st.FrameWidth, sy+st.FrameHeight)).(*ebiten.Image), op)
}
