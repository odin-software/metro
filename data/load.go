package data

import (
	"log"

	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/models"
)

func LoadStations() []*models.Station {
	db := baso.NewBaso()
	stations, err := db.ListStations()
	if err != nil {
		log.Fatal(err)
	}
	result := make([]*models.Station, 0)
	for _, station := range stations {
		result = append(
			result,
			models.NewStation(station.ID, station.Name, station.Position),
		)
	}
	return result
}

func LoadLines(stations []*models.Station) []models.Line {
	db := baso.NewBaso()
	lines := db.ListLinesWithStations()

	// Build station lookup map by ID for O(1) access
	stationsByID := make(map[int64]*models.Station)
	for _, st := range stations {
		stationsByID[st.ID] = st
	}

	result := make([]models.Line, 0)
	for _, line := range lines {
		// Match line's stations to the loaded station pointers
		stationPtrs := make([]*models.Station, 0, len(line.Stations))
		for _, st := range line.Stations {
			if stPtr, ok := stationsByID[st.ID]; ok {
				stationPtrs = append(stationPtrs, stPtr)
			}
		}

		result = append(
			result,
			models.Line{
				ID:       line.ID,
				Name:     line.Name,
				Stations: stationPtrs,
			},
		)
	}

	return result
}

func LoadTrains(
	stations []*models.Station,
	lines []models.Line,
	central *models.Network[models.Station],
	eventChannel chan<- interface{},
) []models.Train {
	db := baso.NewBaso()
	trainsData := db.ListTrainsFull()
	makes := db.ListMakes()

	// Build lookup maps for O(1) access
	makesByName := make(map[string]models.Make)
	for _, make := range makes {
		makesByName[make.Name] = make
	}

	linesByName := make(map[string]models.Line)
	for _, line := range lines {
		linesByName[line.Name] = line
	}

	stationsById := make(map[int64]*models.Station)
	for _, station := range stations {
		stationsById[station.ID] = station
	}

	result := make([]models.Train, 0)
	for _, train := range trainsData {
		mk, ok := makesByName[train.MakeName]
		if !ok {
			log.Fatalf("Make not found: %s", train.MakeName)
		}

		line, ok := linesByName[train.LineName]
		if !ok {
			log.Fatalf("Line not found: %s", train.LineName)
		}

		st, ok := stationsById[train.CurrentStationId]
		if !ok {
			log.Fatalf("Station not found with ID: %d", train.CurrentStationId)
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
				eventChannel,
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
		// Converting the type from the database into the memory model.
		st1 := models.Station{
			ID:       station1.ID,
			Name:     station1.Name,
			Position: models.NewVector(station1.Position.X, station1.Position.Y),
		}
		st2 := models.Station{
			ID:       station2.ID,
			Name:     station2.Name,
			Position: models.NewVector(station2.Position.X, station2.Position.Y),
		}
		cn.InsertEdge(st1, st2, eps)
	}
}
