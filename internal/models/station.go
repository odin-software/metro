package models

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/odin-software/metro/internal/assets"
	"github.com/odin-software/metro/internal/broadcast"
)

type Station struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Position   Vector  `json:"position"`
	Trains     []Train `json:"trains"`
	arrivals   <-chan broadcast.ADMessage[Train]
	departures <-chan broadcast.ADMessage[Train]
	Drawing
}

func NewStation(id int64, name string, location Vector, arr <-chan broadcast.ADMessage[Train], dep <-chan broadcast.ADMessage[Train]) *Station {
	img, frameWidth, frameHeight, frameCount := assets.GetStationSprite()
	st := &Station{
		ID:         id,
		Name:       name,
		Position:   location,
		Trains:     []Train{},
		arrivals:   arr,
		departures: dep,
		Drawing: Drawing{
			Counter:     0,
			FrameWidth:  frameWidth,
			FrameHeight: frameHeight,
			FrameCount:  frameCount,
			Sprite:      img,
		},
	}

	go st.ListenForArrivals()
	go st.ListenForDepartures()

	return st
}

func (st *Station) Update() {
	st.Drawing.Counter++
}

func (st *Station) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(st.FrameWidth)/2, -float64(st.FrameHeight)/2)
	op.GeoM.Translate(st.Position.X, st.Position.Y)
	i := (st.Counter / st.FrameCount) % st.FrameCount
	sx, sy := 0+i*st.FrameWidth, 0
	screen.DrawImage(st.Sprite.SubImage(image.Rect(sx, sy, sx+st.FrameWidth, sy+st.FrameHeight)).(*ebiten.Image), op)
}

func (st *Station) AddTrain(train Train) {
	st.Trains = append(st.Trains, train)
}

func (st *Station) RemoveTrain(train Train) {
	for i, t := range st.Trains {
		if t.Name == train.Name {
			st.Trains = append(st.Trains[:i], st.Trains[i+1:]...)
			break
		}
	}
}

func (st *Station) ListenForArrivals() {
	// This should be a goroutine that listens for trains and adds them to the station.
	for msg := range st.arrivals {
		if msg.StationID != st.ID {
			continue
		}
		st.AddTrain(msg.Train)
	}
}

func (st *Station) ListenForDepartures() {
	// This should be a goroutine that listens for trains and removes them to the station.
	for msg := range st.departures {
		if msg.StationID != st.ID {
			continue
		}
		st.RemoveTrain(msg.Train)
	}
}
