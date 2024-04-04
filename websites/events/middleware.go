package events

import "github.com/labstack/echo/v4"

func Handshake(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
