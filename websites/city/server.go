package city

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/baso"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func Server() {
	server := echo.New()
	bs := baso.NewBaso()
	server.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: control.LoggingFormat,
		},
	))

	server.Static("/ce-js", "websites/city/js")
	server.Static("/ce-css", "websites/city/css")
	server.Static("/ce-images", "websites/city/images")

	server.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, Index())
	})
	server.GET("/editor", func(c echo.Context) error {
		return Render(c, http.StatusOK, Editor())
	})
	server.GET("/stations", func(c echo.Context) error {
		stations := bs.ListStations()
		return c.JSON(http.StatusOK, stations)
	})

	port := fmt.Sprintf(":%d", control.DefaultConfig.PortCity)
	server.Logger.Fatal(server.Start(port))
}
