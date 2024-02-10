package model

import "testing"

func TestNewStation(t *testing.T) {
	station := NewStation("station-1", "Station", NewVector(2.0, 4.0))

	if station.Name != "Station" {
		t.Fatal("The station was wrongly created.")
	}
	if station.Location.X != 2.0 {
		t.Fatal("The station was wrongly created.")
	}
}
