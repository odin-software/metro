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

	// Pages
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ticker.Resume()
		Index().Render(r.Context(), w)
	})
	mux.HandleFunc("/editor", func(w http.ResponseWriter, r *http.Request) {
		ticker.Pause()
		Editor().Render(r.Context(), w)
	})

	// server.GET("/edges", func(c echo.Context) error {
	// 	edges := bs.ListEdges()
	// 	return c.JSON(http.StatusOK, edges)
	// })
	// server.GET("/edges/:id", func(c echo.Context) error {
	// 	stringId := c.Param("id")
	// 	id, err := strconv.Atoi(stringId)
	// 	if err != nil {
	// 		return c.NoContent(400)
	// 	}
	// 	edges := bs.ListEdgePoints(int64(id))
	// 	return c.JSON(http.StatusOK, edges)
	// })

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
