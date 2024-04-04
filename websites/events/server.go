package events

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/odin-software/metro/control"
)

func Server() {
	server := echo.New()

	port := fmt.Sprintf(":%d", control.DefaultConfig.PortCity)
	server.Logger.Fatal(server.Start(port))
}
