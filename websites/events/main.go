package events

import (
	"fmt"
	"net/http"

	"github.com/odin-software/metro/control"
	"golang.org/x/net/websocket"
)

func handleBase() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("hit /")
	})
}

func Main() {
	server := NewServer()
	router := http.NewServeMux()

	router.Handle("GET /trains", websocket.Handler(server.handleTrains))
	router.Handle("GET /", handleBase())

	port := fmt.Sprintf(":%d", control.DefaultConfig.PortEvents)
	http.ListenAndServe(port, router)
}
