package city

import (
	"encoding/json"
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

	err = s.baso.CreateStations(reqStations)
	if err != nil {
		InternalServerErrorHandler(w, req)
		return
	}

	w.WriteHeader(http.StatusOK)
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
