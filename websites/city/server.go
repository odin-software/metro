package city

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/sematick"
)

type Server struct {
	basoMux sync.Mutex
	baso    *baso.Baso
	ticker  *sematick.Ticker
}

type CreateLine struct {
	Name     string `json:"name"`
	Stations []struct {
		StationId int64 `json:"stationId"`
		Odr       int64 `json:"odr"`
	} `json:"stations"`
	Color string `json:"color"`
}

type MoveTrainToLine struct {
	TrainId int64 `json:"trainId"`
	LineId  int64 `json:"lineId"`
}

type GenerateNetworkParams struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	Radius int `json:"radius"`
}

func NewServer(tick *sematick.Ticker) *Server {
	return &Server{
		baso:   baso.NewBaso(),
		ticker: tick,
	}
}

func InternalServerErrorHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func BadRequestErrorHandler(w http.ResponseWriter, req *http.Request, msg string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(msg))
}

func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func JsonHandler(w http.ResponseWriter, req *http.Request, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		InternalServerErrorHandler(w, req)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (s *Server) GetAllStations(w http.ResponseWriter, req *http.Request) {
	s.basoMux.Lock()
	defer s.basoMux.Unlock()

	stations, err := s.baso.ListStations()
	if err != nil {
		InternalServerErrorHandler(w, req)
		return
	}
	if len(stations) == 0 {
		NotFoundHandler(w, req)
		return
	}

	JsonHandler(w, req, stations)
}

func (s *Server) CreateStations(w http.ResponseWriter, req *http.Request) {
	s.basoMux.Lock()
	defer s.basoMux.Unlock()

	var reqStations []baso.CreateStation
	err := json.NewDecoder(req.Body).Decode(&reqStations)
	if err != nil {
		BadRequestErrorHandler(w, req, "Malformed request body.")
		return
	}

	newStations, err := s.baso.CreateStations(reqStations)
	if err != nil {
		InternalServerErrorHandler(w, req)
		return
	}

	JsonHandler(w, req, newStations)
}

func (s *Server) GetLines(w http.ResponseWriter, req *http.Request) {
	s.basoMux.Lock()
	defer s.basoMux.Unlock()

	lines, err := s.baso.ListLinesWithPoints()
	if err != nil {
		InternalServerErrorHandler(w, req)
		return
	}
	if len(lines) == 0 {
		NotFoundHandler(w, req)
		return
	}

	JsonHandler(w, req, lines)
}

func (s *Server) CreateLine(w http.ResponseWriter, req *http.Request) {
	s.basoMux.Lock()
	defer s.basoMux.Unlock()

	var reqLine CreateLine
	err := json.NewDecoder(req.Body).Decode(&reqLine)
	if err != nil {
		BadRequestErrorHandler(w, req, "Malformed request body.")
		return
	}

	id, err := s.baso.CreateLine(reqLine.Name, reqLine.Color)
	if err != nil {
		InternalServerErrorHandler(w, req)
	}

	for _, st := range reqLine.Stations {
		_, err := s.baso.GetStationById(st.StationId)
		if err != nil {
			BadRequestErrorHandler(w, req, "Station sent does not exist.")
		}
		_, err = s.baso.CreateStationLine(st.StationId, id, st.Odr)
		if err != nil {
			InternalServerErrorHandler(w, req)
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetEdges(w http.ResponseWriter, req *http.Request) {
	s.basoMux.Lock()
	defer s.basoMux.Unlock()

	edges, err := s.baso.ListEdges()
	if err != nil {
		InternalServerErrorHandler(w, req)
		return
	}
	if len(edges) == 0 {
		NotFoundHandler(w, req)
		return
	}

	JsonHandler(w, req, edges)
}

func (s *Server) CreateEdges(w http.ResponseWriter, req *http.Request) {
	s.basoMux.Lock()
	defer s.basoMux.Unlock()

	var reqEdges []baso.CreateEdge
	err := json.NewDecoder(req.Body).Decode(&reqEdges)
	if err != nil {
		BadRequestErrorHandler(w, req, "Malformed request body.")
		return
	}

	edges, err := s.baso.CreateEdges(reqEdges)
	if err != nil {
		BadRequestErrorHandler(w, req, "an id from the edge does not exists")
		return
	}

	JsonHandler(w, req, edges)
}

func (s *Server) GetEdgePoints(w http.ResponseWriter, req *http.Request) {
	s.basoMux.Lock()
	defer s.basoMux.Unlock()

	stringId := req.PathValue("id")
	id, err := strconv.Atoi(stringId)
	if err != nil {
		InternalServerErrorHandler(w, req)
		return
	}

	edgePoints, err := s.baso.ListEdgePoints(int64(id))
	if err != nil {
		InternalServerErrorHandler(w, req)
		return
	}
	if len(edgePoints) == 0 {
		NotFoundHandler(w, req)
		return
	}

	JsonHandler(w, req, edgePoints)
}

func (s *Server) GetTrains(w http.ResponseWriter, req *http.Request) {
	s.basoMux.Lock()
	defer s.basoMux.Unlock()

	trains := s.baso.ListTrainsFull()
	if len(trains) == 0 {
		NotFoundHandler(w, req)
		return
	}

	JsonHandler(w, req, trains)
}

func (s *Server) UpdateTrainToLine(w http.ResponseWriter, req *http.Request) {
	s.basoMux.Lock()
	defer s.basoMux.Unlock()

	var reqMove MoveTrainToLine
	err := json.NewDecoder(req.Body).Decode(&reqMove)
	if err != nil {
		BadRequestErrorHandler(w, req, "Malformed request body.")
		return
	}

	err = s.baso.MoveTrainToLine(reqMove.TrainId, reqMove.LineId)
	if err != nil {
		BadRequestErrorHandler(w, req, "an id from the edge does not exists")
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) GenerateNetwork(w http.ResponseWriter, req *http.Request) {
	s.basoMux.Lock()
	defer s.basoMux.Unlock()

	var reqNetworkParams GenerateNetworkParams
	err := json.NewDecoder(req.Body).Decode(&reqNetworkParams)
	if err != nil {
		BadRequestErrorHandler(w, req, "Malformed request body.")
		return
	}

	err = s.baso.WipeData()
	if err != nil {
		InternalServerErrorHandler(w, req)
	}

	vecs := poissonDiskSampling(float64(reqNetworkParams.Radius), reqNetworkParams.Width, reqNetworkParams.Height, 30)
	stationsToCreate := make([]baso.CreateStation, 0)
	for _, v := range vecs {
		name := strconv.FormatFloat(rand.Float64()*100000.00, 'E', -1, 64)
		stationsToCreate = append(stationsToCreate, baso.CreateStation{
			Name:  name,
			X:     v.X,
			Y:     v.Y,
			Z:     0.0,
			Color: "#48FF23",
		})
	}
	newStations, err := s.baso.CreateStations(stationsToCreate)
	if err != nil {
		InternalServerErrorHandler(w, req)
		return
	}

	w.Write([]byte(fmt.Sprintf("%d stations created.", len(newStations))))
}
