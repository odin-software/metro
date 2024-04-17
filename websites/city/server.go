package city

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/sematick"
)

type CreateStationReq struct {
	Name string  `json:"name"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func Server(tick *sematick.Ticker) {
	server := echo.New()
	bs := baso.NewBaso()
	server.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: control.LoggingFormat,
		},
	))

	server.Static("/ce-js", "websites/city/dist")
	server.Static("/ce-css", "websites/city/css")
	server.Static("/ce-images", "websites/city/images")

	server.GET("/", func(c echo.Context) error {
		tick.Resume()
		return Render(c, http.StatusOK, Index())
	})
	server.GET("/editor", func(c echo.Context) error {
		tick.Pause()
		return Render(c, http.StatusOK, Editor())
	})

	server.GET("/stations", func(c echo.Context) error {
		stations := bs.ListStations()
		return c.JSON(http.StatusOK, stations)
	})
	server.GET("/lines", func(c echo.Context) error {
		lines := bs.ListLinesWithPoints()
		return c.JSON(http.StatusOK, lines)
	})
	server.GET("/edges", func(c echo.Context) error {
		edges := bs.ListEdges()
		return c.JSON(http.StatusOK, edges)
	})
	server.GET("/edges/:id", func(c echo.Context) error {
		stringId := c.Param("id")
		id, err := strconv.Atoi(stringId)
		if err != nil {
			return c.NoContent(400)
		}
		edges := bs.ListEdgePoints(int64(id))
		return c.JSON(http.StatusOK, edges)
	})

	server.POST("/stations", func(c echo.Context) error {
		stReq := new([]CreateStationReq)
		if err := c.Bind(stReq); err != nil {
			return err
		}

		for _, r := range *stReq {
			err := bs.CreateStation(r.Name, r.X, r.Y, 0.0)
			if err != nil {
				return err
			}
		}

		return c.NoContent(http.StatusCreated)
	})

	server.GET("/pause", func(c echo.Context) error {
		tick.Pause()
		return c.NoContent(http.StatusOK)
	})
	server.GET("/resume", func(c echo.Context) error {
		tick.Resume()
		return c.NoContent(http.StatusOK)
	})

	port := fmt.Sprintf(":%d", control.DefaultConfig.PortCity)
	server.Logger.Fatal(server.Start(port))
}
