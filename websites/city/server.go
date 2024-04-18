package city

import (
	"net/http"
	"sync"

	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/sematick"
)

type CreateStationReq struct {
	Name string  `json:"name"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
}

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

func NotFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("404 Not Found"))
}

func (s *Server) GetAllStations(w http.ResponseWriter, req *http.Request) {
	s.basoMux.Lock()
	defer s.basoMux.Unlock()

	stations := s.baso.ListStations()

	if len(stations) == 0 {
		NotFoundHandler(w, req)
	}
}
