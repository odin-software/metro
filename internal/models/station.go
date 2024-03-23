package models

import "github.com/odin-software/metro/internal/broadcast"

type Station struct {
	ID         int64
	Name       string
	Position   Vector
	Trains     []Train
	arrivals   <-chan broadcast.ADMessage[Train]
	departures <-chan broadcast.ADMessage[Train]
}

func NewStation(id int64, name string, location Vector, arr <-chan broadcast.ADMessage[Train], dep <-chan broadcast.ADMessage[Train]) *Station {
	st := &Station{
		ID:         id,
		Name:       name,
		Position:   location,
		Trains:     []Train{},
		arrivals:   arr,
		departures: dep,
	}

	go st.ListenForArrivals()
	go st.ListenForDepartures()

	return st
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
