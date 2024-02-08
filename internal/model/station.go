package model

type Station struct {
	ID       string
	name     string
	location Vector
}

func NewStation(id string, name string, location Vector) Station {
	return Station{
		ID:       id,
		name:     name,
		location: location,
	}
}
