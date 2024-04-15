package events

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/odin-software/metro/control"
	"github.com/odin-software/metro/internal/baso"
	"github.com/odin-software/metro/internal/loglytics"
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

func (s *Server) handleTrains(ws *websocket.Conn) {
	s.conmux.Lock()

	fmt.Println(ws.RemoteAddr(), "is now connected to the trains feed.")
	bs := baso.NewBaso()

	s.conns[ws] = true

	s.conmux.Unlock()

	go s.readLoop(ws)

	for {
		trains := bs.ListTrainsFull()
		payload, err := json.Marshal(trains)
		if err != nil {
			fmt.Println(err)
		}
		ws.Write([]byte(payload))
		time.Sleep(control.DefaultConfig.WSTrainDuration)
	}
}

func (s *Server) handleLogs(ws *websocket.Conn) {
	s.conmux.Lock()

	fmt.Println(ws.RemoteAddr(), "is now connected to the logs feed.")

	s.conns[ws] = true
	s.conmux.Unlock()

	go s.readLoop(ws)

	for {
		files := loglytics.GetOrderedLogFiles()
		lines := loglytics.GetLastLines(files[len(files)-1], 6)

		payload, err := json.Marshal(lines)
		if err != nil {
			fmt.Println(err)
		}
		ws.Write([]byte(payload))

		time.Sleep(control.DefaultConfig.WSLogsDuration)
	}
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
