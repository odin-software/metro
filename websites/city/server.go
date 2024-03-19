package City

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func CityEditorServer() {
	server := echo.New()
	server.Use(middleware.Logger())

	server.Static("/ce-js", "websites/city-editor/js")
	server.Static("/ce-css", "websites/city-editor/css")
	server.Static("/ce-images", "websites/city-editor/images")

	server.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, Index())
	})
	server.GET("/editor", func(c echo.Context) error {
		return Render(c, http.StatusOK, Editor())
	})

	server.Logger.Fatal(server.Start(":2221"))
}
