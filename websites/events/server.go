package events

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

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

func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client to orderbook feed:", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(time.Second * 1)
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
		s.broadcast(msg)
		// fmt.Println(string(msg))
		// ws.Write([]byte("thank you for the msg!"))
	}
}

func (s *Server) broadcast(b []byte) {
	s.conmux.Lock()
	defer s.conmux.Unlock()
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error: ", err)
			}
		}(ws)
	}
}

func Main() {
	server := NewServer()
	router := http.NewServeMux()

	router.Handle("GET /wso", websocket.Handler(server.handleWSOrderbook))

	http.Handle("/ws", websocket.Handler(server.handleWS))

	port := fmt.Sprintf(":%d", control.DefaultConfig.PortEvents)
	http.ListenAndServe(port, router)
}
