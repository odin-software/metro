package models

type Line = struct {
	ID       int64
	Name     string
	Stations []Station
}
