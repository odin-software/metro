package data

import (
	"errors"
	"log"

	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/broadcast"
	"github.com/odin-software/metro/internal/models"
)

func LoadStations(arrs broadcast.BroadcastServer[broadcast.ADMessage[models.Train]], deps broadcast.BroadcastServer[broadcast.ADMessage[models.Train]]) []*models.Station {
	db := baso.NewBaso()
	stations, err := db.ListStations()
	if err != nil {
		log.Fatal(err)
	}
	result := make([]*models.Station, 0)
	for _, station := range stations {
		result = append(
			result,
			models.NewStation(station.ID, station.Name, station.Position, arrs.Subscribe(), deps.Subscribe()),
		)
	}
	return result
}

func LoadLines() []models.Line {
	db := baso.NewBaso()
	lines := db.ListLinesWithStations()
	result := make([]models.Line, 0)
	for _, line := range lines {
		result = append(
			result,
			models.Line{
				Name:     line.Name,
				Stations: line.Stations,
			},
		)
	}

	return result
}

func LoadTrains(
	stations []*models.Station,
	lines []models.Line,
	central *models.Network[models.Station],
	a chan broadcast.ADMessage[models.Train],
	d chan broadcast.ADMessage[models.Train],
) []models.Train {
	db := baso.NewBaso()
	trainsData := db.ListTrainsFull()
	makes := db.ListMakes()
	result := make([]models.Train, 0)
	for _, train := range trainsData {
		mk, err := getMakeByName(train.MakeName, makes)
		if err != nil {
			log.Fatal(err)
		}
		line, err := getLineByName(train.LineName, lines)
		if err != nil {
			log.Fatal(err)
		}
		st, err := getStationById(train.CurrentStationId, stations)
		if err != nil {
			log.Fatal(err)
		}
		result = append(
			result,
			models.NewTrain(
				train.Name,
				mk,
				models.NewVector(train.X, train.Y),
				st,
				line,
				central,
				a,
				d,
			),
		)
	}
	return result
}

func LoadEdges(cn *models.Network[models.Station]) {
	db := baso.NewBaso()
	edges, err := db.ListEdges()
	if err != nil {
		log.Fatal(err)
	}
	for _, edge := range edges {
		edgePoints, err := db.ListEdgePoints(edge.ID)
		if err != nil {
			log.Fatal(err)
		}
		station1, err := db.GetStationById(edge.Fromid)
		if err != nil {
			log.Fatal(err)
		}
		station2, err := db.GetStationById(edge.Toid)
		if err != nil {
			log.Fatal(err)
		}
		eps := make([]models.Vector, 0)
		for _, ep := range edgePoints {
			eps = append(eps, models.NewVector(ep.X, ep.Y))
		}
		cn.InsertEdge(station1, station2, eps)
	}
}

func LoadEverything(
	arrs broadcast.BroadcastServer[broadcast.ADMessage[models.Train]],
	deps broadcast.BroadcastServer[broadcast.ADMessage[models.Train]],
	cn *models.Network[models.Station],
) {
	// lines := LoadLines()
	// stations := LoadStations(arrs, deps)
	LoadEdges(cn)
}

func getMakeByName(name string, makes []models.Make) (models.Make, error) {
	for _, make := range makes {
		if make.Name == name {
			return make, nil
		}
	}
	err := errors.New("make not found")
	return models.Make{}, err
}

func getLineByName(name string, lines []models.Line) (models.Line, error) {
	for _, line := range lines {
		if line.Name == name {
			return line, nil
		}
	}
	err := errors.New("line not found")
	return models.Line{}, err
}

func getStationById(id int64, stations []*models.Station) (models.Station, error) {
	for _, station := range stations {
		if station.ID == id {
			return *station, nil
		}
	}
	err := errors.New("station not found")
	return models.Station{}, err
}
