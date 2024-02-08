package model

import (
	"data"
	"fmt"
)

type Make struct {
	name        string
	description string
	accMag      float64
	topSpeed    float64
}

type Train struct {
	name         string
	make         Make
	Position     Vector
	velocity     Vector
	acceleration Vector
	current      Station
	next         Station
	destinations data.Queue[Station]
	q            data.Queue[Station]
	central      *data.Graph[Station]
}

func NewTrain(name string, make Make, pos Vector, central *data.Graph[Station]) Train {
	return Train{
		name:         name,
		make:         make,
		Position:     pos,
		velocity:     NewVector(0.0, 0.0),
		acceleration: NewVector(0.0, 0.0),
		current:      Station{},
		next:         Station{},
		destinations: data.Queue[Station]{},
		q:            data.Queue[Station]{},
		central:      central,
	}
}

func (tr *Train) AddDestination(st Station) {
	tr.destinations.Q(st)
}

func (tr *Train) addToQueue(sts []Station) {
	tr.q.QList(sts)
}

func (tr *Train) Update() {
	// If there is no next station, assign one from the destinations queue
	if tr.next == (Station{}) {
		tr.next, _ = tr.destinations.DQ()
		path, err := tr.central.ShortestPath(tr.current, tr.next)
		if err != nil {
			fmt.Println("Something went wrong")
		}
		tr.addToQueue(path)
	}

	// Update velocity based of direction of next location
	reach, err := tr.q.Peek()
	if err != nil {
		fmt.Println("No items in queue (?)")
		return
	}

	direction := reach.Location.SoftSub(tr.Position)
	mag := direction.Magnitude()
	where := tr.current.Location.Dist(tr.next.Location) / 10

	if mag < where {
		m := Map(mag, 0, where, 0, tr.make.accMag)
		direction.SetMag(m)
	} else {
		direction.SetMag(tr.make.accMag)
	}

	// Update position based on velocity
	tr.velocity.Add(direction)
	tr.velocity.Limit(tr.make.topSpeed)
	tr.Position.Add(tr.velocity)

	if tr.Position.Dist(reach.Location) < 0.02 {
		st, err := tr.q.DQ()
		if err != nil {
			fmt.Println("Something wrong with the next q value.")
		}
		tr.current = st
		if st == tr.next {
			tr.destinations.DQ()
			tr.next = Station{}
		}
	}
}

func NewMake(name string, description string, accMag float64, topSpeed float64) Make {
	return Make{
		name:        name,
		description: description,
		accMag:      accMag,
		topSpeed:    topSpeed,
	}
}
