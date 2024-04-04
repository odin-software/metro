package City

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/odin-software/metro/internal/sematick"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func CityServer(ticker *sematick.Ticker) {
	server := echo.New()
	server.Use(middleware.Logger())

	server.Static("/ce-js", "websites/city/js")
	server.Static("/ce-css", "websites/city/css")
	server.Static("/ce-images", "websites/city/images")

	server.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, Index())
	})
	server.GET("/pause", func(c echo.Context) error {
		ticker.Pause()
		return nil
	})
	server.GET("/resume", func(c echo.Context) error {
		ticker.Resume()
		return nil
	})
	server.GET("/editor", func(c echo.Context) error {
		return Render(c, http.StatusOK, Editor())
	})

	server.Logger.Fatal(server.Start(":2221"))
}
