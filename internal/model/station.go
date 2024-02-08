package model

type Station struct {
	id       string
	name     string
	location Vector
}

func NewStation(id string, name string, location Vector) Station {
	return Station{
		id:       id,
		name:     name,
		location: location,
	}
}
