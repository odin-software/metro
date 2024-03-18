package VirtualWorld

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

func VirtualWorldServer() {
	server := echo.New()
	server.Use(middleware.Logger())

	server.Static("/vw-js", "websites/virtual-world/js")
	server.Static("/vw-css", "websites/virtual-world/css")
	server.Static("/vw-images", "websites/virtual-world/images")

	server.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, Index())
	})

	server.Logger.Fatal(server.Start(":2445"))
}
