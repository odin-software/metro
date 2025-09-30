package models

import (
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/assets"
)

type Make struct {
	Name        string
	Description string
	AccMag      float64
	TopSpeed    float64
}

type Train struct {
	Name         string
	make         Make
	Position     Vector
	velocity     Vector
	acceleration Vector
	Current      Station
	Next         *Station
	forward      bool
	destinations Line
	q            Queue[Vector]
	central      *Network[Station]
	waitCounter  int // Ticks to wait at station (non-blocking)
	Drawing
}

func NewTrain(
	name string,
	make Make,
	pos Vector,
	initialStation Station,
	line Line,
	central *Network[Station],
) Train {
	img, frameWidth, frameHeight, frameCount := assets.GetTrainSprite()
	return Train{
		Name:         name,
		make:         make,
		Position:     pos,
		velocity:     NewVector(0.0, 0.0),
		acceleration: NewVector(0.0, 0.0),
		Current:      initialStation,
		Next:         nil,
		forward:      true,
		destinations: line,
		q:            Queue[Vector]{},
		central:      central,
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
	return Make{
		Name:        name,
		Description: description,
		AccMag:      accMag,
		TopSpeed:    topSpeed,
	}
}

func (tr *Train) addToQueue(sts []Vector) {
	tr.q.QList(sts)
}

func (tr *Train) logArrival(stationName string) {
	logMsg := fmt.Sprintf("%s arrived at station: %s", tr.Name, stationName)
	control.Log(logMsg)
}

func (tr *Train) logDeparture(stationName string) {
	logMsg := fmt.Sprintf("%s departed from station: %s", tr.Name, stationName)
	control.Log(logMsg)
}

func (tr *Train) getNextFromDestinations() *Station {
	var next *Station
	for i, st := range tr.destinations.Stations {
		if st.ID == tr.Current.ID {
			if tr.forward && i == len(tr.destinations.Stations)-1 {
				tr.forward = false
				next = &tr.destinations.Stations[i-1]
				break
			}
			if !tr.forward && i == 0 {
				tr.forward = true
				next = &tr.destinations.Stations[i+1]
				break
			}
			if tr.forward {
				next = &tr.destinations.Stations[i+1]
				break
			}
			next = &tr.destinations.Stations[i-1]
			break
		}
	}

	// This means the train was moved to another line.
	if next == nil {
		next = &tr.destinations.Stations[0]
	}

	return next
}

func (tr *Train) Tick() {
	// If waiting at station, decrement counter and skip this tick
	if tr.waitCounter > 0 {
		tr.waitCounter--
		return
	}

	// If there is no next station, assign one from the destinations queue
	if tr.Next == nil {
		tr.Next = tr.getNextFromDestinations()

		// Adding points between the current station and the next one.
		path, err := tr.central.AreConnected(tr.Current, *tr.Next)
		path = append(path, tr.Next.Position)
		if err != nil {
			control.Log(fmt.Sprintf("Error connecting stations %s to %s: %v", tr.Current.Name, tr.Next.Name, err))
		}
		tr.addToQueue(path)

		tr.logDeparture(tr.Current.Name)
	}

	// Update velocity based of direction of next location
	reach, err := tr.q.Peek()
	if err != nil {
		control.Log(fmt.Sprintf("Train %s: No items in queue", tr.Name))
		return
	}

	direction := reach.SoftSub(tr.Position)
	mag := direction.Magnitude()
	where := tr.Current.Position.Dist(reach) / 8

	// Slow down if we are close.
	if mag < where {
		m := Map(mag, 0, where, 0, tr.make.AccMag)
		direction.SetMag(m)
	} else {
		direction.SetMag(tr.make.AccMag)
	}

	// Update position based on velocity
	tr.velocity.Add(direction)
	tr.velocity.Limit(tr.make.TopSpeed)
	tr.Position.Add(tr.velocity)
	distance := tr.Position.Dist(reach)

	if distance <= 1 {
		_, err := tr.q.DQ()
		if err != nil {
			return
		}
		tr.velocity.Scale(0)
		if tr.q.Size() == 0 {
			tr.Current = *tr.Next
			tr.Next = nil

			// Log arrival
			tr.logArrival(tr.Current.Name)

			// Calculate wait ticks: waitDuration / tickDuration
			waitTicks := int(control.DefaultConfig.TrainWaitInStation / control.DefaultConfig.LoopDuration)
			tr.waitCounter = waitTicks
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
