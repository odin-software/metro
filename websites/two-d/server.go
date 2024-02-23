package TwoD

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

func TwoDimensionalServer() {
	server := echo.New()
	server.Use(middleware.Logger())

	server.Static("/web/2d/images", "images")
	server.Static("/web/2d/css", "css")

	server.GET("/", func(c echo.Context) error {
		return Render(c, http.StatusOK, hello("Ken"))
	})

	server.Logger.Fatal(server.Start(":2444"))
}
