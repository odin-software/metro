package model

import (
	"fmt"
	"internal/broadcast"
)

type Station struct {
	ID         string
	Name       string
	Position   Vector
	Trains     []Train
	arrivals   <-chan broadcast.ADMessage[Train]
	departures <-chan broadcast.ADMessage[Train]
}

func NewStation(id string, name string, location Vector, arr <-chan broadcast.ADMessage[Train], dep <-chan broadcast.ADMessage[Train]) Station {
	st := Station{
		ID:         id,
		Name:       name,
		Position:   location,
		Trains:     []Train{},
		arrivals:   arr,
		departures: dep,
	}

	go st.ListenForTrains()

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

func (st *Station) ListenForTrains() {
	// This should be a goroutine that listens for trains and adds them to the station.
	select {
	case msg := <-st.arrivals:
		fmt.Println("Arrival")
		if msg.StationID != st.ID {
			return
		}
		st.AddTrain(msg.Train)
	case msg := <-st.departures:
		fmt.Println("Departures")
		if msg.StationID != st.ID {
			return
		}
		st.RemoveTrain(msg.Train)
	default:
		fmt.Println("No message")
	}
}
