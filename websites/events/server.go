package events

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/odin-software/metro/control"
)

const WSHeaderKey = "Sec-WebSocket-Key"
const MagicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func Server() {
	server := echo.New()

	server.GET("/", func(c echo.Context) error {
		if c.Request().Header.Get(WSHeaderKey) != control.DefaultConfig.WebSocketKey {
			return echo.NewHTTPError(400)
		}
		c.Response().Header().Set("Upgrade", "websocket")
		return c.HTML(http.StatusSwitchingProtocols, "")
	})

	port := fmt.Sprintf(":%d", control.DefaultConfig.PortEvents)
	server.Logger.Fatal(server.Start(port))
}
