package model

type Station struct {
	ID       string
	Name     string
	Location Vector
}

func NewStation(id string, name string, location Vector) Station {
	return Station{
		ID:       id,
		Name:     name,
		Location: location,
	}
}
