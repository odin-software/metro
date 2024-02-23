package model

import (
	"fmt"
	"internal/broadcast"
	"time"
)

type Make struct {
	name        string
	description string
	accMag      float64
	topSpeed    float64
}

type Train struct {
	Name         string
	make         Make
	Position     Vector
	velocity     Vector
	acceleration Vector
	current      Station
	next         *Station
	forward      bool
	destinations Line
	q            Queue[Vector]
	central      *Network[Station]
	arrivals     chan<- broadcast.ADMessage[Train]
	departures   chan<- broadcast.ADMessage[Train]
}

func NewTrain(
	name string,
	make Make,
	pos Vector,
	initialStation Station,
	line Line,
	central *Network[Station],
	a chan broadcast.ADMessage[Train],
	d chan broadcast.ADMessage[Train],
) Train {
	return Train{
		Name:         name,
		make:         make,
		Position:     pos,
		velocity:     NewVector(0.0, 0.0),
		acceleration: NewVector(0.0, 0.0),
		current:      initialStation,
		next:         nil,
		forward:      true,
		destinations: line,
		q:            Queue[Vector]{},
		central:      central,
		arrivals:     a,
		departures:   d,
	}
}

func (tr *Train) addToQueue(sts []Vector) {
	tr.q.QList(sts)
}

func (tr *Train) Update() {
	// If there is no next station, assign one from the destinations queue
	if tr.next == nil {
		if tr.forward {
			for i, st := range tr.destinations.Stations {
				if st.ID == tr.current.ID {
					if i == len(tr.destinations.Stations)-1 {
						tr.forward = false
						tr.next = &tr.destinations.Stations[i-1]
						break
					}
					tr.next = &tr.destinations.Stations[i+1]
					break
				}
			}
		} else {
			for i, st := range tr.destinations.Stations {
				if st.ID == tr.current.ID {
					if i == 0 {
						tr.forward = true
						tr.next = &tr.destinations.Stations[i+1]
						break
					}
					tr.next = &tr.destinations.Stations[i-1]
					break
				}
			}
		}
		path, err := tr.central.AreConnected(tr.current, *tr.next)
		path = append(path, tr.next.Position)
		if err != nil {
			fmt.Println("Something went wrong")
		}
		tr.addToQueue(path)

		// Calculate time to the next station
		// This needs to be updated to check for arrival.
		// Maybe later I could come up with an equation.
		pc := tr.Position.Copy()
		vc := NewVector(0, 0)
		ticks := 0
		for i := 0; i < len(path); i++ {
			for pc.Dist(path[i]) >= 1 {
				ticks++
				d := path[i].SoftSub(pc)
				d.SetMag(tr.make.accMag)
				vc.Add(d)
				vc.Limit(tr.make.topSpeed)
				pc.Add(vc)
			}
			vc.Scale(0)
		}

		timeToNext := float64(ticks) / 60.0
		fmt.Printf("%s is going to %s, it will arrive in %.1f seconds\n", tr.Name, tr.next.Name, timeToNext)

		// Broadcast departure
		msg := broadcast.ADMessage[Train]{
			StationID: tr.current.ID,
			Train:     *tr,
		}
		tr.departures <- msg
	}

	// Update velocity based of direction of next location
	reach, err := tr.q.Peek()
	if err != nil {
		fmt.Println("No items in queue (?)")
		return
	}

	direction := reach.SoftSub(tr.Position)
	mag := direction.Magnitude()
	where := tr.current.Position.Dist(reach) / 8

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
	distance := tr.Position.Dist(reach)

	if distance <= 1 {
		tr.q.DQ()
		tr.velocity.Scale(0)
		if tr.q.Size() == 0 {
			tr.current = *tr.next
			tr.next = nil
			fmt.Printf("%s arrived at: %s\n", tr.Name, tr.current.Name)

			// Broadcast arrival
			msg := broadcast.ADMessage[Train]{
				StationID: tr.current.ID,
				Train:     *tr,
			}
			tr.arrivals <- msg

			time.Sleep(3 * time.Second)
			return
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
