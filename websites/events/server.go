package events

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/odin-software/metro/control"
	"golang.org/x/net/websocket"
)

type Server struct {
	conmux sync.Mutex
	conns  map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	s.conmux.Lock()
	defer s.conmux.Unlock()
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}

		msg := buf[:n]
		fmt.Println(string(msg))
		ws.Write([]byte("thank you for the msg!"))
	}
}

const WSHeaderKey = "Sec-WebSocket-Key"
const MagicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func Main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	port := fmt.Sprintf(":%d", control.DefaultConfig.PortEvents)
	http.ListenAndServe(port, nil)
}
