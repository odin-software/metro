package events

import (
	"fmt"
	"net/http"

	"github.com/odin-software/metro/control"
	"golang.org/x/net/websocket"
)

func Main() {
	server := NewServer()
	router := http.NewServeMux()

	router.Handle("GET /trains", websocket.Handler(server.handleTrains))

	port := fmt.Sprintf(":%d", control.DefaultConfig.PortEvents)
	http.ListenAndServe(port, router)
}
