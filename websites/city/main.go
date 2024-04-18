package city

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/sematick"
)

func Main(tick *sematick.Ticker) {
	mux := http.NewServeMux()
	server := NewServer(tick)

	// Static directories
	jsFs := http.FileServer(http.Dir("websites/city/dist"))
	cssFs := http.FileServer(http.Dir("websites/city/css"))
	mux.Handle("/js/", http.StripPrefix("/js/", jsFs))
	mux.Handle("/css/", http.StripPrefix("/css/", cssFs))

	// Pages
	mux.Handle("/", templ.Handler(Index()))
	mux.Handle("GET /editor", templ.Handler(Editor()))

	// Stations endpoints
	mux.HandleFunc("GET /stations", server.GetAllStations)

	// Stations endpoints
	// server.GET("/", func(c echo.Context) error {
	// 	tick.Resume()
	// 	return Render(c, http.StatusOK, Index())
	// })
	// server.GET("/editor", func(c echo.Context) error {
	// 	tick.Pause()
	// 	return Render(c, http.StatusOK, Editor())
	// })

	// server.GET("/stations", func(c echo.Context) error {
	// 	stations := bs.ListStations()
	// 	return c.JSON(http.StatusOK, stations)
	// })
	// server.GET("/lines", func(c echo.Context) error {
	// 	lines := bs.ListLinesWithPoints()
	// 	return c.JSON(http.StatusOK, lines)
	// })
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

	// server.POST("/stations", func(c echo.Context) error {
	// 	stReq := new([]CreateStationReq)
	// 	if err := c.Bind(stReq); err != nil {
	// 		return err
	// 	}

	// 	for _, r := range *stReq {
	// 		err := bs.CreateStation(r.Name, r.X, r.Y, 0.0)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}

	// 	return c.NoContent(http.StatusCreated)
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
