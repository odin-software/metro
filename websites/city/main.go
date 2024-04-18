package city

import (
	"fmt"
	"log"
	"net/http"

	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/sematick"
)

func Main(ticker *sematick.Ticker) {
	mux := http.NewServeMux()
	server := NewServer(ticker)

	// Static directories
	jsFs := http.FileServer(http.Dir("websites/city/dist"))
	cssFs := http.FileServer(http.Dir("websites/city/css"))
	mux.Handle("/js/", http.StripPrefix("/js/", jsFs))
	mux.Handle("/css/", http.StripPrefix("/css/", cssFs))

	// Stations endpoints
	mux.HandleFunc("GET /stations", server.GetAllStations)
	mux.HandleFunc("POST /stations", server.CreateStations)

	// Lines endpoints
	mux.HandleFunc("GET /lines", server.GetLines)
	mux.HandleFunc("GET /edges", server.GetEdges)
	mux.HandleFunc("GET /edges/{id}", server.GetEdgePoints)

	// Ticker endpoints
	mux.HandleFunc("GET /pause", func(w http.ResponseWriter, r *http.Request) {
		ticker.Pause()
		w.WriteHeader(http.StatusNoContent)
	})

	// Pages
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ticker.Resume()
		Index().Render(r.Context(), w)
	})
	mux.HandleFunc("/editor", func(w http.ResponseWriter, r *http.Request) {
		ticker.Pause()
		Editor().Render(r.Context(), w)
	})

	// server.GET("/pause", func(c echo.Context) error {
	// 	tick.Pause()
	// 	return c.NoContent(http.StatusOK)
	// })
	// server.GET("/resume", func(c echo.Context) error {
	// 	tick.Resume()
	// 	return c.NoContent(http.StatusOK)
	// })

	port := fmt.Sprintf(":%d", control.DefaultConfig.PortCity)
	log.Printf("server is listening at %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}
