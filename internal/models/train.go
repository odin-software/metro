package models

import (
	"fmt"
	"image"
	_ "image/png"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/assets"
)

type Make struct {
	Name         string
	Description  string
	AccMag       float64 // Acceleration in pixels/tick²
	TopSpeed     float64 // Top speed in pixels/tick
	TopSpeedKmH  float64 // Top speed in km/h (real-world)
	AccelerationMPS2 float64 // Acceleration in m/s² (real-world)
}

// EventEmitter is an interface for emitting train events to Tenjin
type EventEmitter interface {
	Emit(event interface{})
}

type Train struct {
	Name           string
	model          Make // Renamed from 'make' to avoid conflict with built-in
	Position       Vector
	velocity       Vector
	Current        *Station // Pointer to avoid copying mutex
	Next           *Station
	forward        bool
	destinations   Line
	q              Queue[Vector]
	central        *Network[Station]
	waitCounter    int                // Ticks to wait at station (non-blocking)
	waitTicks      int                // Precomputed wait duration in ticks
	eventChannel   chan<- interface{} // Channel to send events to Tenjin
	tickCounter    int                // Counter for periodic tick events (emit every 60 ticks)
	Capacity       int                // Maximum number of passengers
	Passengers     []*Passenger       // Current passengers on board
	passengerMutex sync.RWMutex       // Thread safety for passenger operations
	Drawing
}

func NewTrain(
	name string,
	trainMake Make,
	pos Vector,
	initialStation *Station,
	line Line,
	central *Network[Station],
	eventChannel chan<- interface{},
) Train {
	img, frameWidth, frameHeight, frameCount := assets.GetTrainSprite()
	// Precompute wait duration in ticks
	waitTicks := int(control.DefaultConfig.TrainWaitInStation / control.DefaultConfig.LoopDuration)
	return Train{
		Name:         name,
		model:        trainMake,
		Position:     pos,
		velocity:     NewVector(0.0, 0.0),
		Current:      initialStation,
		Next:         nil,
		forward:      true,
		destinations: line,
		q:            Queue[Vector]{},
		central:      central,
		waitTicks:    waitTicks,
		eventChannel: eventChannel,
		tickCounter:  0,
		Capacity:     50, // Default capacity: 50 passengers
		Passengers:   make([]*Passenger, 0),
		Drawing: Drawing{
			Counter:     0,
			FrameWidth:  frameWidth,
			FrameHeight: frameHeight,
			FrameCount:  frameCount,
			Sprite:      img,
		},
	}
}

func NewMake(name string, description string, accMag float64, topSpeed float64) Make {
	// Calculate real-world speeds
	topSpeedKmH := PixelSpeedToKmPerHour(topSpeed)

	// Calculate acceleration in m/s²
	// accMag is in pixels/tick², convert to m/s²
	ticksPerSecond := 1.0 / control.DefaultConfig.LoopDuration.Seconds()
	pixelsPerSecondSquared := accMag * ticksPerSecond * ticksPerSecond
	metersPerSecondSquared := PixelsToMeters(pixelsPerSecondSquared)

	return Make{
		Name:              name,
		Description:       description,
		AccMag:            accMag,
		TopSpeed:          topSpeed,
		TopSpeedKmH:       topSpeedKmH,
		AccelerationMPS2:  metersPerSecondSquared,
	}
}

func (tr *Train) addToQueue(sts []Vector) {
	tr.q.QList(sts)
}

func (tr *Train) logArrival(stationName string) {
	// Emit arrival event to Tenjin
	if tr.eventChannel != nil && control.DefaultConfig.TenjinEnabled {
		event := struct {
			Type        string
			Train       string
			StationID   int64
			StationName string
			Time        time.Time
			Position    Vector
		}{
			Type:        "train_arrival",
			Train:       tr.Name,
			StationID:   tr.Current.ID,
			StationName: stationName,
			Time:        time.Now(),
			Position:    tr.Position,
		}
		select {
		case tr.eventChannel <- event:
		default:
			// Channel full, skip event (non-blocking)
		}
	}
}

func (tr *Train) logDeparture(stationName string) {
	logMsg := fmt.Sprintf("%s departed from station: %s", tr.Name, stationName)
	control.Log(logMsg)

	// Emit departure event to Tenjin
	if tr.eventChannel != nil && control.DefaultConfig.TenjinEnabled {
		nextName := ""
		if tr.Next != nil {
			nextName = tr.Next.Name
		}
		event := struct {
			Type        string
			Train       string
			StationID   int64
			StationName string
			NextStation string
			Time        time.Time
			Position    Vector
		}{
			Type:        "train_departure",
			Train:       tr.Name,
			StationID:   tr.Current.ID,
			StationName: stationName,
			NextStation: nextName,
			Time:        time.Now(),
			Position:    tr.Position,
		}
		select {
		case tr.eventChannel <- event:
		default:
			// Channel full, skip event (non-blocking)
		}
	}
}

func (tr *Train) getNextFromDestinations() *Station {
	var next *Station
	for i, st := range tr.destinations.Stations {
		if st.ID == tr.Current.ID {
			if tr.forward && i == len(tr.destinations.Stations)-1 {
				tr.forward = false
				next = tr.destinations.Stations[i-1]
				break
			}
			if !tr.forward && i == 0 {
				tr.forward = true
				next = tr.destinations.Stations[i+1]
				break
			}
			if tr.forward {
				next = tr.destinations.Stations[i+1]
				break
			}
			next = tr.destinations.Stations[i-1]
			break
		}
	}

	// This means the train was moved to another line.
	if next == nil {
		next = tr.destinations.Stations[0]
	}

	return next
}

func (tr *Train) Tick() {
	// Increment tick counter for periodic events
	tr.tickCounter++

	// Emit tick event every 60 ticks (once per second)
	if tr.tickCounter >= 60 {
		tr.emitTickEvent()
		tr.tickCounter = 0
	}

	// If waiting at station, decrement counter and skip this tick
	if tr.waitCounter > 0 {
		tr.waitCounter--
		return
	}

	// If there is no next station, assign one from the destinations queue
	if tr.Next == nil {
		tr.Next = tr.getNextFromDestinations()

		// Adding points between the current station and the next one.
		path, err := tr.central.AreConnected(*tr.Current, *tr.Next)
		path = append(path, tr.Next.Position)
		if err != nil {
			errMsg := fmt.Sprintf("Error connecting stations %s to %s: %v", tr.Current.Name, tr.Next.Name, err)
			control.Log(errMsg)
			tr.emitErrorEvent(errMsg, "path_connection")
		}
		tr.addToQueue(path)

		tr.logDeparture(tr.Current.Name)
	}

	// Update velocity based of direction of next location
	reach, err := tr.q.Peek()
	if err != nil {
		errMsg := fmt.Sprintf("Train %s: No items in queue", tr.Name)
		control.Log(errMsg)
		tr.emitErrorEvent(errMsg, "empty_queue")
		return
	}

	direction := reach.SoftSub(tr.Position)
	mag := direction.Magnitude()
	where := tr.Current.Position.Dist(reach) / 4 // Larger slowing zone for faster trains

	// Slow down if we are close - scale velocity directly for visible deceleration.
	if mag < where {
		m := Map(mag, 0, where, 0, tr.model.AccMag)
		direction.SetMagFrom(mag, m) // Reuse calculated magnitude
	} else {
		direction.SetMagFrom(mag, tr.model.AccMag) // Reuse calculated magnitude
	}

	// Update position based on velocity
	tr.velocity.Add(direction)
	tr.velocity.Limit(tr.model.TopSpeed)
	tr.Position.Add(tr.velocity)
	distance := tr.Position.Dist(reach)

	if distance <= 1 {
		_, err := tr.q.DQ()
		if err != nil {
			return
		}
		tr.velocity.Scale(0)
		if tr.q.Size() == 0 {
			tr.Current = tr.Next
			tr.Next = nil

			// Log arrival
			tr.logArrival(tr.Current.Name)

			// Passenger operations
			tr.handlePassengerDisembark()
			tr.handlePassengerBoarding()

			// Use precomputed wait ticks
			tr.waitCounter = tr.waitTicks
			return
		}
	}
}

// Display methods

func (tr *Train) Update() {
	tr.Drawing.Counter++
}

func (tr *Train) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(tr.FrameWidth)/2, -float64(tr.FrameHeight)/2)
	op.GeoM.Translate(tr.Position.X, tr.Position.Y)
	i := (tr.Counter / tr.FrameCount) % tr.FrameCount
	sx, sy := 0+i*tr.FrameWidth, 0
	screen.DrawImage(tr.Sprite.SubImage(image.Rect(sx, sy, sx+tr.FrameWidth, sy+tr.FrameHeight)).(*ebiten.Image), op)
}

// emitTickEvent sends periodic state updates to Tenjin
func (tr *Train) emitTickEvent() {
	if tr.eventChannel == nil || !control.DefaultConfig.TenjinEnabled {
		return
	}

	nextStationID := int64(0)
	if tr.Next != nil {
		nextStationID = tr.Next.ID
	}

	event := struct {
		Type           string
		Train          string
		Position       Vector
		Velocity       Vector
		Speed          float64
		CurrentStation int64
		NextStation    int64
		Time           time.Time
	}{
		Type:           "train_tick",
		Train:          tr.Name,
		Position:       tr.Position,
		Velocity:       tr.velocity,
		Speed:          tr.velocity.Magnitude(),
		CurrentStation: tr.Current.ID,
		NextStation:    nextStationID,
		Time:           time.Now(),
	}

	select {
	case tr.eventChannel <- event:
	default:
		// Channel full, skip event (non-blocking)
	}
}

// emitErrorEvent sends error events to Tenjin
func (tr *Train) emitErrorEvent(errMsg, context string) {
	if tr.eventChannel == nil || !control.DefaultConfig.TenjinEnabled {
		return
	}

	event := struct {
		Type    string
		Train   string
		Error   string
		Context string
		Time    time.Time
	}{
		Type:    "train_error",
		Train:   tr.Name,
		Error:   errMsg,
		Context: context,
		Time:    time.Now(),
	}

	select {
	case tr.eventChannel <- event:
	default:
		// Channel full, skip event (non-blocking)
	}
}

// Passenger management methods

// AddPassenger adds a passenger to the train
func (tr *Train) AddPassenger(passenger *Passenger) bool {
	tr.passengerMutex.Lock()
	defer tr.passengerMutex.Unlock()

	// Check if train is full
	if len(tr.Passengers) >= tr.Capacity {
		return false
	}

	tr.Passengers = append(tr.Passengers, passenger)
	return true
}

// RemovePassenger removes a passenger from the train
func (tr *Train) RemovePassenger(passenger *Passenger) bool {
	tr.passengerMutex.Lock()
	defer tr.passengerMutex.Unlock()

	for i, p := range tr.Passengers {
		if p.ID == passenger.ID {
			// Remove passenger by swapping with last and truncating
			tr.Passengers[i] = tr.Passengers[len(tr.Passengers)-1]
			tr.Passengers = tr.Passengers[:len(tr.Passengers)-1]
			return true
		}
	}
	return false
}

// GetPassengerCount returns the current number of passengers
func (tr *Train) GetPassengerCount() int {
	tr.passengerMutex.RLock()
	defer tr.passengerMutex.RUnlock()
	return len(tr.Passengers)
}

// IsFull returns true if train is at capacity
func (tr *Train) IsFull() bool {
	tr.passengerMutex.RLock()
	defer tr.passengerMutex.RUnlock()
	return len(tr.Passengers) >= tr.Capacity
}

// IsCrowded returns true if train is over 80% capacity
func (tr *Train) IsCrowded() bool {
	tr.passengerMutex.RLock()
	defer tr.passengerMutex.RUnlock()
	return float64(len(tr.Passengers))/float64(tr.Capacity) > 0.8
}

// GetCapacityPercentage returns the percentage of capacity used
func (tr *Train) GetCapacityPercentage() float64 {
	tr.passengerMutex.RLock()
	defer tr.passengerMutex.RUnlock()
	if tr.Capacity == 0 {
		return 0
	}
	return (float64(len(tr.Passengers)) / float64(tr.Capacity)) * 100
}

// GetSpeed returns the current speed (magnitude of velocity in pixels/tick)
func (tr *Train) GetSpeed() float64 {
	return tr.velocity.Magnitude()
}

// GetSpeedKmH returns the current speed in km/h
func (tr *Train) GetSpeedKmH() float64 {
	return PixelSpeedToKmPerHour(tr.velocity.Magnitude())
}

// GetDistanceToNext returns the distance to the next waypoint in meters
func (tr *Train) GetDistanceToNext() float64 {
	reach, err := tr.q.Peek()
	if err != nil {
		return 0
	}
	pixelDistance := tr.Position.Dist(reach)
	return PixelsToMeters(pixelDistance)
}

// GetPassengers returns a copy of the passengers slice (thread-safe)
func (tr *Train) GetPassengers() []*Passenger {
	tr.passengerMutex.RLock()
	defer tr.passengerMutex.RUnlock()

	passengers := make([]*Passenger, len(tr.Passengers))
	copy(passengers, tr.Passengers)
	return passengers
}

// GetPassengersForStation returns passengers destined for a specific station
func (tr *Train) GetPassengersForStation(stationID int64) []*Passenger {
	tr.passengerMutex.RLock()
	defer tr.passengerMutex.RUnlock()

	result := make([]*Passenger, 0)
	for _, p := range tr.Passengers {
		if p.DestinationStation.ID == stationID {
			result = append(result, p)
		}
	}
	return result
}

// handlePassengerDisembark removes passengers who have reached their destination
func (tr *Train) handlePassengerDisembark() {
	if tr.Current == nil {
		return
	}

	passengers := tr.GetPassengersForStation(tr.Current.ID)
	for _, p := range passengers {
		tr.RemovePassenger(p)
		p.DisembarkTrain(tr.Current)
	}
}

// handlePassengerBoarding boards waiting passengers up to capacity
func (tr *Train) handlePassengerBoarding() {
	if tr.Current == nil || tr.IsFull() {
		return
	}

	// Get waiting passengers at this station
	waiting := tr.Current.GetWaitingPassengers()
	for _, p := range waiting {
		if tr.IsFull() {
			break
		}
		// Board passenger
		tr.Current.RemovePassenger(p)
		tr.AddPassenger(p)
		p.BoardTrain(tr)
	}
}
